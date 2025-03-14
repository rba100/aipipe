using System;
using System.Text;
using System.Text.RegularExpressions;
using System.Linq;
using System.Collections.Generic;

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
    private static readonly Regex numberedListRegex = new(@"^(\s*)(\d+\.)\s+(.*)$", RegexOptions.Compiled);
    private static readonly Regex unorderedListRegex = new(@"^(\s*)([-*])\s+(.*)$", RegexOptions.Compiled);
    private static readonly Regex emphasisRegex = new(@"(\*\*|__)(.*?)\1|(\*|_)(.*?)\3", RegexOptions.Compiled);
    private static readonly Regex blockQuoteRegex = new(@"^(\s*)((?:>\s*)+)(.*)$", RegexOptions.Compiled);
    private static readonly Regex horizontalRuleRegex = new(@"^(\s*)([-*_])\2\2+\s*$", RegexOptions.Compiled);

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

        if (horizontalRuleRegex.IsMatch(line))
        {
            PrintHorizontalRule(line);
            return;
        }

        if (blockQuoteRegex.IsMatch(line))
        {
            PrintBlockQuote(line);
            return;
        }

        if (numberedListRegex.IsMatch(line))
        {
            PrintNumberedList(line);
            return;
        }

        if (unorderedListRegex.IsMatch(line))
        {
            PrintUnorderedList(line);
            return;
        }

        PrintFormattedText(line);
    }

    private void PrintFormattedText(string line)
    {
        var lastIndex = 0;
        var inlineCodeMatches = inlineCodeRegex.Matches(line);
        var emphasisMatches = emphasisRegex.Matches(line);
        
        // Combine and sort all matches by index
        var allMatches = new List<(int Index, int Length, string Type)>();
        
        foreach (Match match in inlineCodeMatches)
        {
            allMatches.Add((match.Index, match.Length, "code"));
        }
        
        foreach (Match match in emphasisMatches)
        {
            allMatches.Add((match.Index, match.Length, "emphasis"));
        }
        
        allMatches = allMatches.OrderBy(m => m.Index).ToList();
        
        foreach (var (index, length, type) in allMatches)
        {
            // Print text before the match
            if (index > lastIndex)
            {
                SetColor(ConsoleColor.Green);
                Console.Write(line[lastIndex..index]);
            }
            
            // Print the match with appropriate formatting
            if (type == "code")
            {
                SetColor(ConsoleColor.Cyan);
                Console.Write(line[index..(index + length)]);
            }
            else if (type == "emphasis")
            {
                SetColor(ConsoleColor.White);
                if (isBoldSupported)
                {
                    Console.Write("\u001b[1m");
                    Console.Write(line[index..(index + length)]);
                    Console.Write("\u001b[22m");
                }
                else
                {
                    Console.Write(line[index..(index + length)]);
                }
            }
            
            lastIndex = index + length;
        }
        
        // Print remaining text
        if (lastIndex < line.Length)
        {
            SetColor(ConsoleColor.Green);
            Console.Write(line[lastIndex..]);
        }
    }

    private void PrintNumberedList(string line)
    {
        var match = numberedListRegex.Match(line);
        if (match.Success)
        {
            var indentation = match.Groups[1].Value;
            var number = match.Groups[2].Value;
            var content = match.Groups[3].Value;

            // Print indentation
            Console.Write(indentation);

            // Print number with highlight
            SetColor(ConsoleColor.Blue);
            Console.Write(number);

            // Print content with formatting
            Console.Write(" ");
            PrintFormattedText(content);
        }
    }

    private void PrintUnorderedList(string line)
    {
        var match = unorderedListRegex.Match(line);
        if (match.Success)
        {
            var indentation = match.Groups[1].Value;
            var bullet = match.Groups[2].Value;
            var content = match.Groups[3].Value;

            // Print indentation
            Console.Write(indentation);

            // Print bullet with highlight
            SetColor(ConsoleColor.Blue);
            Console.Write(bullet);

            // Print content with formatting
            Console.Write(" ");
            PrintFormattedText(content);
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

    private void PrintBlockQuote(string line)
    {
        var match = blockQuoteRegex.Match(line);
        if (match.Success)
        {
            var indentation = match.Groups[1].Value;
            var quoteMarkers = match.Groups[2].Value;
            var content = match.Groups[3].Value;

            // Print indentation
            Console.Write(indentation);

            // Extract and print each '>' character in blue
            var markers = quoteMarkers.TrimEnd();
            for (int i = 0; i < markers.Length; i++)
            {
                if (markers[i] == '>')
                {
                    SetColor(ConsoleColor.Blue);
                    Console.Write('>');
                }
                else
                {
                    Console.Write(markers[i]);
                }
            }

            // Print content in cyan
            SetColor(ConsoleColor.Cyan);
            Console.Write(content);
        }
    }

    private void PrintHorizontalRule(string line)
    {
        SetColor(ConsoleColor.Yellow);
        Console.Write(line);
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
