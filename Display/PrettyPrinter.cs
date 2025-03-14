using System;
using System.Text;
using System.Text.RegularExpressions;

namespace aipipe.Display;

public enum PrintState
{
    Normal,
    InCodeBlock
}

public class PrettyPrinter : IDisposable
{
    private readonly ConsoleColor originalColor;
    private readonly bool isBoldSupported;
    private readonly StringBuilder lineBuffer = new();
    private PrintState currentState = PrintState.Normal;
    
    private static readonly Regex headerRegex = new(@"^#{1,6}\s+.*$", RegexOptions.Compiled);
    private static readonly Regex inlineCodeRegex = new(@"`[^`\n]+`", RegexOptions.Compiled);
    private static readonly Regex codeBlockStartRegex = new(@"^```", RegexOptions.Compiled);
    private static readonly Regex codeBlockEndRegex = new(@"^```\s*$", RegexOptions.Compiled);

    public PrettyPrinter()
    {
        originalColor = Console.ForegroundColor;
        isBoldSupported = !OperatingSystem.IsWindows() || Environment.GetEnvironmentVariable("WT_SESSION") != null;
    }

    public void Print(string text)
    {
        if(text.Length == 0)
        {
            return;
        }
        if (!text.Contains('\n'))
        {
            lineBuffer.Append(text);
            return;
        }

        var isTerminated = text.EndsWith('\n');
        text = isTerminated ? text[..^1] : text;

        var lines = text.Split('\n');
        for (int i = 0; i < lines.Length; i++)
        {
            var line = lines[i];
            bool isLastLine = i == lines.Length - 1;

            if (lineBuffer.Length > 0)
            {
                line = lineBuffer.ToString() + line;
                lineBuffer.Clear();
            }

            if (isLastLine && !isTerminated)
            {
                lineBuffer.Append(line);
                return;
            }

            ProcessLine(line);
            if (!isLastLine)
                Console.WriteLine();
        }
        if (isTerminated)
            Console.WriteLine();
    }

    private void ProcessLine(string line)
    {
        if(line.Contains("\r")) throw new ArgumentException("Line should not contain carriage return characters");
        if(line.Contains("\n")) throw new ArgumentException("Line should not contain newline characters");
        if (currentState == PrintState.Normal)
        {
            if (codeBlockStartRegex.IsMatch(line))
            {
                SetColor(ConsoleColor.Cyan);
                Console.Write(line);
                currentState = PrintState.InCodeBlock;
                return;
            }

            ProcessNormalLine(line);
        }
        else // InCodeBlock
        {
            if (codeBlockEndRegex.IsMatch(line))
            {
                SetColor(ConsoleColor.Cyan);
                Console.Write(line);
                currentState = PrintState.Normal;
                return;
            }

            SetColor(ConsoleColor.Cyan);
            Console.Write(line);
        }
    }

    private void ProcessNormalLine(string line)
    {
        if (headerRegex.IsMatch(line))
        {
            PrintHeader(line);
            return;
        }

        var lastIndex = 0;
        var inlineMatches = inlineCodeRegex.Matches(line);

        foreach (Match match in inlineMatches)
        {
            // Print text before inline code
            if (match.Index > lastIndex)
            {
                SetColor(ConsoleColor.Green);
                Console.Write(line[lastIndex..match.Index]);
            }

            // Print inline code
            SetColor(ConsoleColor.Magenta);
            Console.Write(match.Value);
            lastIndex = match.Index + match.Length;
        }

        // Print remaining text
        if (lastIndex < line.Length)
        {
            SetColor(ConsoleColor.Green);
            Console.Write(line[lastIndex..]);
        }
    }

    private void PrintHeader(string header)
    {
        SetColor(ConsoleColor.Yellow);
        
        if (isBoldSupported)
        {
            Console.Write("\u001b[1m");
            Console.Write(header);
            Console.Write("\u001b[22m");
        }
        else
        {
            var hashCount = header.TakeWhile(c => c == '#').Count();
            var headerText = header[hashCount..].TrimStart();
            Console.Write(new string('#', hashCount) + " " + headerText.ToUpperInvariant());
        }
    }

    private void SetColor(ConsoleColor color)
    {
        if (Console.ForegroundColor != color)
            Console.ForegroundColor = color;
    }

    public void Dispose()
    {
        if (lineBuffer.Length > 0)
        {
            ProcessLine(lineBuffer.ToString());
            Console.WriteLine();
        }
        
        SetColor(originalColor);
    }
}
