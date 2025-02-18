using System.Threading.Tasks;

namespace aipipe.llms;

public interface ILLMClient
{
    Task<string> CreateCompletionAsync(string prompt);
}
