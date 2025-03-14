using System.CommandLine;

namespace aipipe;

public class CommandLineOptions
{
    public Option<bool> CodeBlockOption { get; set; } = new Option<bool>(aliases: new[] { "--cb", "-c" }, description: "Extract code block from response");
    public Option<bool> ReasoningOption { get; set; } = new Option<bool>(aliases: new[] { "--reasoning", "-r" }, description: "Use a reasoning model");
    public Option<bool> FastOption { get; set; } = new Option<bool>(aliases: new[] { "--fast", "-f" }, description: "Use the Fast model");
    public Option<bool> MicOption { get; set; } = new Option<bool>(aliases: new[] { "--mic", "-m" }, description: "Use microphone input");
    public Option<bool> StreamOption { get; set; } = new Option<bool>(aliases: new[] { "--stream", "-s" }, description: "Stream completions from the AI model");
    public Argument<string?> PromptArgument { get; set; } = new Argument<string?>(name:"prompt", description: "The prompt to send to the AI. Optional, but you must supply at least one input to the AI (prompt, --mic, or pipe in a file)", getDefaultValue: () => null);
    public Option<bool> OpenRouterOption { get; set; } = new Option<bool>(aliases: new[] { "--or", "-o" }, description: "Use OpenRouter");
    public Option<bool> PrettyOption { get; set; } = new Option<bool>(
        aliases: new[] { "--pretty", "-p" },
        description: "Enable pretty printing with colors and formatting"
    );

    public RootCommand RootCommand { get; set; }

    public CommandLineOptions()
    {
        RootCommand = new RootCommand("aipipe - A tool to pipe input to an AI model");
        
        // Add all options and arguments to the command
        RootCommand.AddOption(CodeBlockOption);
        RootCommand.AddOption(ReasoningOption);
        RootCommand.AddOption(FastOption);
        RootCommand.AddOption(MicOption);
        RootCommand.AddArgument(PromptArgument);
        RootCommand.AddOption(OpenRouterOption);
        RootCommand.AddOption(StreamOption);
        RootCommand.AddOption(PrettyOption);
        
        // Add a single validator to the root command to check for mutually exclusive options
        RootCommand.AddValidator(result =>
        {
            if (result.GetValueForOption(CodeBlockOption) && 
                result.GetValueForOption(PrettyOption))
            {
                result.ErrorMessage = "The --cb and --pretty options cannot be used together.";
            }
        });
    }
}
