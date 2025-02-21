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
    private StringBuilder buffer = new();
    private CodeBlockState state = CodeBlockState.SearchingOpening;
    private static readonly Regex openingRegex = new Regex("```[^\n]*\n", RegexOptions.Compiled);

    public string? Handle(string part)
    {
        if (state == CodeBlockState.Closed)
            return null;

        buffer.Append(part);

        if (state == CodeBlockState.SearchingOpening)
        {
            // Check if buffer contains a complete opening delimiter
            var match = openingRegex.Match(buffer.ToString());
            if (match.Success)
            {
                state = CodeBlockState.Open;
                // Remove everything up to the end of the opening delimiter
                buffer.Remove(0, match.Index + match.Length);
            }
            else {
                // Opening delimiter not yet found; do not output any text
                return "";
            }
        }

        string output = "";
        if (state == CodeBlockState.Open)
        {
            string bufStr = buffer.ToString();
            int closePos = bufStr.IndexOf("```");
            if (closePos >= 0)
            {
                // Capture content up to the closing delimiter
                output = bufStr.Substring(0, closePos);
                // Transition to closed state and clear the buffer
                state = CodeBlockState.Closed;
                buffer.Clear();
            }
            else
            {
                // No closing delimiter found, output all and clear buffer
                output = bufStr;
                buffer.Clear();
            }
        }

        return output;
    }
}
