# Chat interface sequence

How `chatModel` in `chat_interface.go` handles a session, from start-up to quitting. Every method in `chat_interface.go` appears below, along with the rendering functions in `rendering.go` and the reply functions in `reply.go` that it calls.

```mermaid
sequenceDiagram
    actor User
    participant RT as Bubble Tea runtime
    participant CM as chatModel<br/>(chat_interface.go)
    participant OM as optionMenu<br/>(chat_interface.go)
    participant R as rendering.go
    participant RP as reply.go
    participant P as providers.Client

    rect rgba(128,128,128,0.08)
    Note over RT,CM: 1. Start-up (main.run)
    RT->>CM: newChatModel(client)
    RT->>CM: Init()
    CM-->>RT: textinput.Blink
    RT->>CM: Update(tea.WindowSizeMsg)
    Note over CM: width = msg.Width
    CM->>CM: updateInput(msg)
    RT->>CM: View()
    CM->>R: renderPromptArea(input.View(), width)
    R-->>RT: prompt with lines above and below
    end

    rect rgba(128,128,128,0.08)
    Note over User,CM: 2. Typing
    User->>RT: key "h"
    RT->>CM: Update(tea.KeyMsg)
    CM->>CM: handleKey(key)
    CM->>CM: updateInput(key)
    RT->>CM: View()
    CM->>R: renderPromptArea(...)
    end

    rect rgba(128,128,128,0.08)
    Note over User,R: 3. Options menu (":" on empty prompt, then up/down/esc)
    User->>RT: key ":"
    RT->>CM: Update(tea.KeyMsg)
    CM->>CM: handleKey(key)
    CM->>OM: navigate(":")
    OM-->>CM: menu.open = true
    RT->>CM: View()
    CM->>R: renderPromptArea(...)
    CM->>R: renderOptions(menu.cursor)
    User->>RT: key "down" / "up" / "esc"
    RT->>CM: Update → handleKey
    CM->>OM: navigate(name)
    OM-->>CM: cursor moved / menu closed
    end

    rect rgba(128,128,128,0.08)
    Note over User,P: 4. Submitting a prompt
    User->>RT: key "enter"
    RT->>CM: Update(tea.KeyMsg)
    CM->>CM: handleKey(key)
    CM->>CM: submit()
    Note over CM: input.Reset(), thinking = true<br/>turn = len(history)/2 + 1
    CM->>R: renderSentPrompt(prompt, turn, width)
    R->>R: renderTurnRule(width), renderLedger(turn, prompt)
    CM->>RP: send(client, history, prompt)
    CM-->>RT: tea.Sequence(Println(sent prompt), send Cmd)
    RT->>User: prints rule and numbered prompt in the ledger
    RT->>CM: View()
    CM->>R: renderPromptArea("Thinking...", width)
    RT->>RP: runs send Cmd in a goroutine
    RP->>P: Send(ctx, history + prompt)
    P-->>RP: Reply, err
    RP-->>RT: replyMsg{prompt, reply, err}
    end

    rect rgba(128,128,128,0.08)
    Note over RT,RP: 5. Reply arrives
    RT->>CM: Update(replyMsg)
    Note over CM: thinking = false<br/>if err == nil: history += prompt, reply.Text
    CM->>R: renderReply(msg)
    R->>RP: replyText(msg)
    RP-->>R: Markdown + usage line, or error
    R->>R: renderLedger("", text)
    R-->>CM: reply in the ledger
    CM-->>RT: tea.Println(reply)
    RT->>User: prints reply above the live view
    RT->>CM: View()
    CM->>R: renderPromptArea(input.View(), width)
    end

    rect rgba(128,128,128,0.08)
    Note over User,CM: 6. Quitting (any of three ways)
    alt ctrl+c / ctrl+d
        RT->>CM: Update → handleKey
        CM-->>RT: tea.Quit
    else menu "Quit" + enter
        RT->>CM: Update → handleKey
        Note over CM: quitting = true
        CM-->>RT: tea.Quit
    else typed "quit" + enter
        RT->>CM: Update → handleKey → submit()
        Note over CM: quitting = true
        CM-->>RT: tea.Quit
    end
    RT->>CM: View()
    CM-->>RT: "Bye..." (unless ctrl+c/d)
    end
```

## Where each function appears

| Function | File:line | Called by | Phase |
|---|---|---|---|
| `newChatModel` | `chat_interface.go:23` | `run` in `main.go` | 1 |
| `Init` | `chat_interface.go:31` | runtime, once | 1 |
| `Update` | `chat_interface.go:37` | runtime, for every message | 1–6 |
| `handleKey` | `chat_interface.go:55` | `Update` on `tea.KeyMsg` | 2, 3, 4, 6 |
| `submit` | `chat_interface.go:72` | `handleKey` on enter | 4, 6 |
| `updateInput` | `chat_interface.go:86` | `Update` (other messages) and `handleKey` (typing) | 1, 2 |
| `View` | `chat_interface.go:93` | runtime, after every `Update` | 1–6 |
| `optionMenu.navigate` | `chat_interface.go:117` | `handleKey` while the menu is open or on ":" | 3 |
| `renderSentPrompt` | `rendering.go:12` | `submit` | 4 |
| `renderReply` | `rendering.go:19` | `Update` on `replyMsg` | 5 |
| `renderPromptArea` | `rendering.go:26` | `View` | 1–5 |
| `renderOptions` | `rendering.go:32` | `View` while the menu is open | 3 |
| `renderLedger` | `rendering.go:49` | `renderSentPrompt`, `renderReply` | 4, 5 |
| `renderTurnRule` | `rendering.go:60` | `renderSentPrompt` | 4 |
| `send` | `reply.go:24` | `submit` | 4 |
| `replyText` | `reply.go:35` | `renderReply` | 5 |

## Notes

- Ctrl+C and Ctrl+D do not set `quitting`, so the last `View` shows the prompt rather than "Bye...". The other two quit paths do set it.
- The menu's enter only closes the list. The `TODO` in `navigate` (`chat_interface.go:126`) is still open, so Provider, Model and History do nothing yet; only the Quit case in `handleKey` acts on a selection.
- The turn number comes from the saved history, which only grows on a successful reply, so a failed send leaves the next prompt with the same number.
- Line numbers reflect the code at the time of writing and will drift as the files change.
