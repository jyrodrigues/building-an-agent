


class ToolUse {
    type "ToolUse"
    toolName "please_read_for_me"
    filename string
}

class TextualReply {
    type "Text"
    value string
}

enum Role {
  Assistant @alias("assistant")
  User @alias("user")
  Tool @alias("tool-response")
}

class Message {
  role Role
  msg string
}

function CallGPT5Mini(chatHistory: Message[]) -> ToolUse | TextualReply {
  // Specify a client as provider/model-name
  // You can also use custom LLM params with a custom client name from clients.baml like "client CustomGPT5" or "client CustomSonnet4"
  client "openai-responses/gpt-5-mini" // Set OPENAI_API_KEY to use this client.
  prompt #"
    {{ _.role("system") }}

    You'll now act as a personalized coding assistant called Red I'll use for didactic purposes.

    I'm teaching how to create AI Agents and how tool use is structured.

    So when asked about the contents of a file you can

    {{ ctx.output_format(
        prefix="Answer in JSON using this schema to access a tool to read a file\n"
      , or_splitter="\nor use reply with a normal message following this JSON schema: \n"
      )
    }}

    {% for m in chatHistory %}
      {# {{ _.role(m.role) }} #}
{{ m.role }}: {{ m.msg }}
    {% endfor %}
  "#
}



// Test the function with a sample resume. Open the VSCode playground to run this.
test my_tool_call {
  functions [CallOpenAiAPI_GPT5Mini]
  args {
    chatHistory [
      {
        role User
        msg "Hello"
      }
      {
        role Assistant
        msg "Hello! I'm Red, your personalized coding assistant. I'm here to help you with coding questions and demonstrate how AI agents use tools."
      }
      {
        role User
        msg "What's in secret.txt?"
      }
      {
        role Assistant
        msg "{ type: 'ToolUse', toolName: 'please_read_for_me', filename: 'secret.txt' }"
      }
      {
        role Tool
        msg "Banana."
      }
      {
        role Assistant
        msg "I found a Banana in secret.txt"
      }
      {
        role User
        msg "What's in another-secret.key?"
      }
    ]
  }
}
