import os

import requests
import streamlit as st
from dotenv import load_dotenv

load_dotenv()

BACKEND_URL = os.environ["BACKEND_URL"]

st.set_page_config(page_title="go-chatbot")
st.title("go-chatbot")

if "user_id" not in st.session_state:
    st.session_state.user_id = None

if not st.session_state.user_id:
    st.caption(
        "Use an existing user ID to resume your previous conversations, "
        "or enter a new one to start fresh."
    )
    user_id = st.text_input("Enter your user ID to continue")
    if user_id:
        st.session_state.user_id = user_id
        st.rerun()
    st.stop()

if "conversation_id" not in st.session_state:
    st.session_state.conversation_id = None
if "messages" not in st.session_state:
    st.session_state.messages = []

headers = {"X-User-ID": st.session_state.user_id}

with st.sidebar:
    st.caption(f"User: {st.session_state.user_id}")
    if st.button("Switch user"):
        st.session_state.user_id = None
        st.session_state.conversation_id = None
        st.session_state.messages = []
        st.rerun()

    st.header("Conversations")
    if st.button("New chat"):
        st.session_state.conversation_id = None
        st.session_state.messages = []
        st.rerun()

    resp = requests.get(f"{BACKEND_URL}/conversations", headers=headers)
    for conv in resp.json() or []:
        if st.button(conv["title"] or "(untitled)", key=conv["id"]):
            detail = requests.get(
                f"{BACKEND_URL}/conversations/{conv['id']}", headers=headers
            ).json()
            st.session_state.conversation_id = conv["id"]
            st.session_state.messages = detail["messages"]
            st.rerun()

for msg in st.session_state.messages:
    with st.chat_message(msg["role"]):
        st.write(msg["content"])

if query := st.chat_input("Message"):
    st.session_state.messages.append({"role": "user", "content": query})
    with st.chat_message("user"):
        st.write(query)

    with st.chat_message("assistant"):
        placeholder = st.empty()
        reply = ""
        payload = {"query": query}
        if st.session_state.conversation_id:
            payload["conversation_id"] = st.session_state.conversation_id

        with requests.post(
            f"{BACKEND_URL}/chat", json=payload, headers=headers, stream=True
        ) as r:
            event = None
            data_lines = []
            done = False
            for line in r.iter_lines():
                if line is None:
                    continue
                if line == b"":
                    data = "\n".join(data_lines)
                    data_lines = []
                    if event == "conversation":
                        st.session_state.conversation_id = data.strip()
                    elif event == "message":
                        reply += data
                        placeholder.write(reply)
                    elif event == "done":
                        done = True
                    elif event == "error":
                        st.error(data)
                        done = True
                    if done:
                        break
                    continue
                line = line.decode("utf-8")
                if line.startswith("event:"):
                    event = line.split(":", 1)[1].strip()
                elif line.startswith("data:"):
                    data_lines.append(line.split(":", 1)[1])

    st.session_state.messages.append({"role": "assistant", "content": reply})
