using NAudio.Wave;
using System.Net.Http.Headers;
using System.Text.Json;

namespace aipipe;

public class Mic
{
    private readonly Config _config;

    public Mic(Config config)
    {
        _config = config;
    }

    public async Task<string?> GetMicInput()
    {
        // Use NAudio to capture audio from the default microphone
        var waveIn = new WaveInEvent();
        waveIn.WaveFormat = new WaveFormat(rate: 16000, bits: 16, channels: 1); // Mono, 16kHz, 16-bit
        
        var recording = new System.IO.MemoryStream();
        var writer = new WaveFileWriter(recording, waveIn.WaveFormat);

        waveIn.DataAvailable += (sender, e) =>
        {
            writer.Write(e.Buffer, 0, e.BytesRecorded);
        };

        waveIn.RecordingStopped += (sender, e) =>
        {
            writer.Flush();
            writer.Dispose();
            waveIn.Dispose();
        };

        waveIn.StartRecording();

        Console.Error.WriteLine("Recording... Press enter to accept or any other key to abort.");
        var key = Console.ReadKey();

        if (key.Key != ConsoleKey.Enter)
        {
            return null;
        }

        waveIn.StopRecording();

        recording.Seek(0, System.IO.SeekOrigin.Begin);

        // Convert audio to text using Whisper API (or any other STT service)
        // This is a placeholder, replace with actual implementation
        var audioData = recording.ToArray();
        var text = await ConvertAudioToText(audioData);

        return text;
    }

    private async Task<string?> ConvertAudioToText(byte[] audioData)
    {
        if (string.IsNullOrEmpty(_config.GroqEndpoint))
        {
            Console.Error.WriteLine("GROQ_ENDPOINT environment variable not set.");
            return null;
        }
        if (string.IsNullOrEmpty(_config.GroqToken))
        {
            Console.Error.WriteLine("GROQ_API_KEY environment variable not set.");
            return null;
        }

        var client = new HttpClient();
        client.BaseAddress = new Uri(_config.GroqEndpoint);
        client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", _config.GroqToken);

        var requestContent = new MultipartFormDataContent();
        var audioContent = new ByteArrayContent(audioData);
        audioContent.Headers.ContentType = MediaTypeHeaderValue.Parse("audio/wav"); // Adjust if using a different format

        requestContent.Add(audioContent, "file", "recording.wav"); // Adjust filename if needed
        requestContent.Add(new StringContent("whisper-large-v3"), "model");

        try
        {
            var response = await client.PostAsync("/openai/v1/audio/transcriptions", requestContent);
            response.EnsureSuccessStatusCode(); // Throw exception if not a success

            var responseString = await response.Content.ReadAsStringAsync();
            
            // Parse the JSON response
            using (JsonDocument document = JsonDocument.Parse(responseString))
            {
                JsonElement root = document.RootElement;
                if (root.TryGetProperty("text", out JsonElement textElement))
                {
                    return textElement.GetString();
                }
                else
                {
                    Console.Error.WriteLine($"Error: 'text' field not found in response: {responseString}");
                    return null;
                }
            }
        }
        catch (HttpRequestException ex)
        {
            Console.Error.WriteLine($"Error during transcription: {ex.Message}");
            return null;
        }
        catch (JsonException ex)
        {
            Console.Error.WriteLine($"Error parsing JSON response: {ex.Message}");
            return null;
        }
    }
}
