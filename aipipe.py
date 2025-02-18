import os
import argparse
import openai
import anthropic
import re
import sys

GROQ_API_KEY = os.getenv('GROQ_API_KEY')
GROQ_ENDPOINT = os.getenv('GROQ_ENDPOINT')
GROQ_MODEL = os.getenv('GROQ_MODEL')

ANTHROPIC_API_KEY = os.getenv('ANTHROPIC_API_KEY')
ANTHROPIC_MODEL = 'claude-3-haiku-20240307'

systemMessage = 'You are a helpful assistant. If the user merely asked a question, do not use a code block. If the user has asked for something written, put it in a code block (```).'

def getGroqCompletion(userMessage, groqModel):
    client = openai.Client(base_url=GROQ_ENDPOINT, api_key=GROQ_API_KEY)
    response = client.chat.completions.create(
        model=groqModel,
        max_tokens=4000,
        messages=[{'role': 'system', 'content': systemMessage},
                  {'role': 'user', 'content': userMessage}]
    )
    return response.choices[0].message.content

def getGpt4Completion(userMessage):
    client = openai.Client()
    response = client.chat.completions.create(
        model="gpt-4-0125-preview",
        messages=[{'role': 'system', 'content': systemMessage},
                  {'role': 'user', 'content': userMessage}]
    )
    return response.choices[0].message.content

def getAnthropicCompletion(userMessage):
    client = anthropic.Anthropic(api_key=ANTHROPIC_API_KEY)
    response = client.messages.create(
        model=ANTHROPIC_MODEL,
        system=systemMessage,
        max_tokens=4000,
        messages=[{'role': 'user', 'content': userMessage}]     
    )
    return response.content[0].text

def extract_code_block(completion):
    code_block_pattern = r'```[a-zA-Z0-9.]*\n([\s\S]+?)\n```'
    match = re.search(code_block_pattern, completion)
    return match.group(1) if match else completion

def main():
    if len(sys.argv) < 2:
        print('Usage: python aifile.py "query" > output.txt')
        return

    parser = argparse.ArgumentParser(description="Get completions from different models.")
    parser.add_argument("prompt", help="The prompt to generate the completion for.")
    parser.add_argument("--codeblock", "--cb", action="store_true", help="Return only the code block in the completion.")
    parser.add_argument("--haiku", action="store_true", help="Use Anthropic Claude model.")
    parser.add_argument("--mx", action="store_true", help="Use Mixtral 8x7b-32768 model.")
    parser.add_argument("--l370", action="store_true", help="Use llama 3 70b.")
    parser.add_argument("--gpt4", action="store_true", help="Use GPT-4 model.")

    args = parser.parse_args()

    useHaiku = args.haiku
    useMixtral = args.mx
    useLlama = args.l370
    useGpt4 = args.gpt4
    codeBlock = args.codeblock
    isatty = sys.stdin.isatty()

    if not isatty:
        filePrompt = sys.stdin.read()

    if args.prompt:
        if(isatty):
            prompt = args.prompt
        else:
            prompt = filePrompt + "\n----\n" + args.prompt
    else:
        prompt = filePrompt

    if not prompt or prompt == '' or prompt.isspace():
        print('Error: No prompt provided', file=sys.stderr)
        return 1

    completion = useHaiku and getAnthropicCompletion(prompt) \
                 or useGpt4 and getGpt4Completion(prompt) \
                 or getGroqCompletion(prompt, useMixtral and 'mixtral-8x7b-32768' or useLlama and 'llama3-70b-8192' or GROQ_MODEL)
    if codeBlock:
        output = extract_code_block(completion)
    else:
        output = completion

    sys.stdout.write(output)

if __name__ == '__main__':
    main()
