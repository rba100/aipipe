using aipipe.Speech;
using NAudio.Wave;

namespace aipipe;

#pragma warning disable CS0162

public class SpeechToText
{
    private readonly Config _config;
    private readonly TranscriptionClient _transcriptionClient;

    // Allow these to be configurable if needed
    private readonly float _silenceThreshold;
    private readonly int _silenceDurationMs;
    private readonly int _minimumRecordingMs;

    private const float DEFAULT_SILENCE_THRESHOLD = 0.015f;
    private const int DEFAULT_SILENCE_DURATION_MS = 2000;
    private const int DEFAULT_MINIMUM_RECORDING_MS = 1000;

    private DateTime _lastDebugOutput = DateTime.MinValue;

    private const bool DEBUG = false;
    private const int DEBUG_OUTPUT_INTERVAL_MS = 100;

    public SpeechToText(Config config, float? silenceThreshold = null, int? silenceDurationMs = null, int? minimumRecordingMs = null)
    {
        _config = config;
        _transcriptionClient = new TranscriptionClient(config);
        _silenceThreshold = silenceThreshold ?? DEFAULT_SILENCE_THRESHOLD;
        _silenceDurationMs = silenceDurationMs ?? DEFAULT_SILENCE_DURATION_MS;
        _minimumRecordingMs = minimumRecordingMs ?? DEFAULT_MINIMUM_RECORDING_MS;
    }

    private static void DebugPrint(string message)
    {
        if (DEBUG)
        {
            Console.Error.WriteLine(message);
        }
    }

    public async Task<string?> GetMicInput(bool useKeyboardInput = true)
    {
        Console.Error.WriteLine($"Starting recording with useKeyboardInput={useKeyboardInput}");
        var waveIn = new WaveInEvent();
        waveIn.WaveFormat = new WaveFormat(rate: 16000, bits: 16, channels: 1);
        Console.Error.WriteLine($"Initialized WaveIn with format: {waveIn.WaveFormat}");

        var recording = new System.IO.MemoryStream();
        var writer = new WaveFileWriter(recording, waveIn.WaveFormat);

        if (useKeyboardInput)
        {
            var taskCompletionSource = new TaskCompletionSource<bool>();

            waveIn.DataAvailable += (sender, e) =>
            {
                writer.Write(e.Buffer, 0, e.BytesRecorded);
            };

            waveIn.StartRecording();

            Console.Error.WriteLine("Recording... Press enter to accept or any other key to abort.");

            try
            {
                var key = Console.ReadKey();
                if (key.Key != ConsoleKey.Enter)
                {
                    return null;
                }
            }
            catch (InvalidOperationException)
            {
                Console.Error.WriteLine("Cannot use keyboard input when console is redirected. Use --auto-detect-silence instead.");
                Environment.Exit(1);
                return null;
            }

            waveIn.StopRecording();
        }
        else
        {
            var silenceStartTime = DateTime.MinValue;
            var recordingStartTime = DateTime.Now;
            var isSilent = true;
            var hasSpokenOnce = false;
            if (DEBUG) {
                Console.Error.WriteLine("Initializing silence detection mode:");
                Console.Error.WriteLine($"  Silence threshold: {_silenceThreshold}");
                Console.Error.WriteLine($"  Silence duration: {_silenceDurationMs}ms");
                Console.Error.WriteLine($"  Minimum recording: {_minimumRecordingMs}ms");
            }

            waveIn.DataAvailable += (sender, e) =>
            {
                writer.Write(e.Buffer, 0, e.BytesRecorded);
                
                var rms = CalculateRMS(e.Buffer);
                var currentTime = DateTime.Now;

                // Throttle debug output
                if (DEBUG && (currentTime - _lastDebugOutput).TotalMilliseconds >= DEBUG_OUTPUT_INTERVAL_MS)
                {
                    Console.Error.WriteLine($"Audio level: {rms:F4} | Silent: {isSilent} | Has spoken: {hasSpokenOnce}");
                    if (isSilent && hasSpokenOnce)
                    {
                        var silenceDuration = (currentTime - silenceStartTime).TotalMilliseconds;
                        var totalDuration = (currentTime - recordingStartTime).TotalMilliseconds;
                        Console.Error.WriteLine($"  Silence duration: {silenceDuration:F0}ms | Total duration: {totalDuration:F0}ms");
                    }
                    _lastDebugOutput = currentTime;
                }
                
                if (rms < _silenceThreshold)
                {
                    if (!isSilent)
                    {
                        if(DEBUG) Console.Error.WriteLine($"\nSilence started at {currentTime:HH:mm:ss.fff}");
                        silenceStartTime = currentTime;
                        isSilent = true;
                    }
                    else if (hasSpokenOnce && 
                            (currentTime - silenceStartTime).TotalMilliseconds > _silenceDurationMs &&
                            (currentTime - recordingStartTime).TotalMilliseconds > _minimumRecordingMs)
                    {
                        if(DEBUG) Console.Error.WriteLine("\nStopping recording due to silence threshold reached:");
                        if(DEBUG) Console.Error.WriteLine($"  Total duration: {(currentTime - recordingStartTime).TotalMilliseconds:F0}ms");
                        if(DEBUG) Console.Error.WriteLine($"  Final silence duration: {(currentTime - silenceStartTime).TotalMilliseconds:F0}ms");
                        waveIn.StopRecording();
                    }
                }
                else
                {
                    if (DEBUG && isSilent)
                    {
                        Console.Error.WriteLine($"\nSpeech detected at {currentTime:HH:mm:ss.fff} (Level: {rms:F4})");
                    }
                    isSilent = false;
                    hasSpokenOnce = true;
                }
            };

            Console.Error.WriteLine("\nStarting recording - Will automatically stop after silence is detected");
            Console.Error.WriteLine("Waiting for speech...\n");
            waveIn.StartRecording();

            // Use simple boolean flag since WaveInEvent will trigger RecordingStopped event
            var isRecording = true;
            waveIn.RecordingStopped += (s, e) => 
            {
                isRecording = false;
                Console.Error.WriteLine("Recording stopped");
            };
            
            while (isRecording)
            {
                await Task.Delay(100);
            }
        }

        Console.Error.WriteLine("\nProcessing recorded audio...");
        writer.Flush();
        
        try 
        {
            recording.Position = 0;
            var audioData = recording.ToArray();
            if (DEBUG) {
                Console.Error.WriteLine($"Captured {audioData.Length} bytes of audio");
                Console.Error.WriteLine("Converting audio to text...");
            }
            
            var text = await ConvertAudioToText(audioData);
            Console.Error.WriteLine(text ?? "Failed to convert audio to text");
            return text;
        }
        finally
        {
            // Clean up all resources in finally block
            writer.Dispose();
            recording.Dispose();
            waveIn.Dispose();
        }
    }

    private float CalculateRMS(byte[] buffer)
    {
        // Convert byte array to 16-bit samples and calculate RMS
        float sum = 0;
        for (int i = 0; i < buffer.Length; i += 2)
        {
            short sample = (short)((buffer[i + 1] << 8) | buffer[i]);
            float normalized = sample / 32768f;
            sum += normalized * normalized;
        }
        return (float)Math.Sqrt(sum / (buffer.Length / 2));
    }

    private async Task<string?> ConvertAudioToText(byte[] audioData)
    {
        return await _transcriptionClient.ConvertAudioToText(audioData);
    }
}
