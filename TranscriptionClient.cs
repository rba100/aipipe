using System.Net.Http.Headers;
using System.Text.Json;

namespace aipipe;

public class TranscriptionClient
{
    private readonly Config _config;

    public TranscriptionClient(Config config)
    {
        _config = config;
    }

    public async Task<string?> ConvertAudioToText(byte[] audioData)
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
        audioContent.Headers.ContentType = MediaTypeHeaderValue.Parse("audio/wav");

        requestContent.Add(audioContent, "file", "recording.wav");
        requestContent.Add(new StringContent("whisper-large-v3"), "model");

        try
        {
            var response = await client.PostAsync("/openai/v1/audio/transcriptions", requestContent);
            response.EnsureSuccessStatusCode();

            var responseString = await response.Content.ReadAsStringAsync();

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