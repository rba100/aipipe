namespace aipipe.llms;

public static class Prompts
{
    public static string GetSystemPrompt(bool isCodeBlock)
    {
        return isCodeBlock
            ? "You are a helpful assistant. If the user has asked for something written, put it in a single code block (```type\\n...\\n```), otherwise just provide the answer. If you do use a codeblock, all other text is ignored."
            : "You are a helpful assistant.";
    }
}
