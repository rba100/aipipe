using System.Text;
using System.Text.RegularExpressions;

namespace aipipe;

enum CodeBlockState
{
    SearchingOpening,
    Open,
    Closed
}

class CodeBlockStreamHandler
{
    private StringBuilder Buffer = new();
    private CodeBlockState state = CodeBlockState.SearchingOpening;
    private static readonly Regex openingRegex = new Regex("```[^\n]*\n", RegexOptions.Compiled);
    private static readonly Regex potentialClosingRegex = new Regex("\n`{0,2}$", RegexOptions.Compiled);

    private IAsyncEnumerable<string> innerStream;
    private bool isCalled = false;

    public CodeBlockStreamHandler(IAsyncEnumerable<string> innerStream)
    {
        this.innerStream = innerStream;
    }

    public async IAsyncEnumerable<string> Stream()
    {
        if (isCalled)
            throw new System.InvalidOperationException("Stream can only be called once");
        isCalled = true;
        await foreach (var part in innerStream)
        {
            var result = Handle(part);
            if (result == null)
                yield break;
            if (result.Length > 0)
                yield return result;
        }
        
        if (state != CodeBlockState.Closed && Buffer.Length > 0)
        {
            yield return Buffer.ToString();
        }
    }

    private string? Handle(string part)
    {
        if (state == CodeBlockState.Closed)
            return null;
        
        Buffer.Append(part);
        string bufStr = Buffer.ToString();

        if (state == CodeBlockState.SearchingOpening)
        {
            var match = openingRegex.Match(bufStr);
            if (match.Success)
            {
                state = CodeBlockState.Open;
                string remainingContent = bufStr.Substring(match.Index + match.Length);
                Buffer.Clear();
                Buffer.Append(remainingContent);
                return Handle(""); // Process remaining content in Open state
            }
            return "";
        }

        if (state == CodeBlockState.Open)
        {
            // Check for potential closing marker first
            if (potentialClosingRegex.IsMatch(bufStr))
            {
                return "";
            }

            int closePos = bufStr.IndexOf("\n```");
            if (closePos >= 0)
            {
                string output = bufStr.Substring(0, closePos);
                state = CodeBlockState.Closed;
                Buffer.Clear();
                return output;
            }

            string output2 = bufStr;
            Buffer.Clear();
            return output2;
        }

        return "";
    }
}
