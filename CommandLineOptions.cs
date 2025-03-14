using System.CommandLine;

namespace aipipe;

public class CommandLineOptions
{
    public Option<bool> CodeBlockOption { get; set; } = new Option<bool>("--cb", "Extract code block from response");
    public Option<bool> ReasoningOption { get; set; } = new Option<bool>("--r", "Use a reasoning model");
    public Option<bool> FastOption { get; set; } = new Option<bool>("--fast", "Use the Fast model");
    public Option<bool> MicOption { get; set; } = new Option<bool>("--mic", "Use microphone input");
    public Option<bool> StreamOption { get; set; } = new Option<bool>("--stream", "Stream completions from the AI model");
    public Argument<string?> PromptArgument { get; set; } = new Argument<string?>(name:"prompt", description: "The prompt to send to the AI. Optional, but you must supply at least one input to the AI (prompt, --mic, or pipe in a file)", getDefaultValue: () => null);
    public Option<bool> OpenRouterOption { get; set; } = new Option<bool>("--or", "Use OpenRouter");
    public Option<bool> PrettyOption { get; set; } = new Option<bool>(
        aliases: new[] { "--pretty", "-p" },
        description: "Enable pretty printing with colors and formatting"
    );

    public RootCommand RootCommand { get; set; }

    public CommandLineOptions()
    {
        RootCommand = new RootCommand("aipipe - A tool to pipe input to an AI model");
        RootCommand.AddOption(CodeBlockOption);
        RootCommand.AddOption(ReasoningOption);
        RootCommand.AddOption(FastOption);
        RootCommand.AddOption(MicOption);
        RootCommand.AddArgument(PromptArgument);
        RootCommand.AddOption(OpenRouterOption);
        RootCommand.AddOption(StreamOption);
        RootCommand.AddOption(PrettyOption);
    }
}
