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
            st.session_state.messages = detail.get("messages") or []
            st.rerun()

for msg in st.session_state.messages:
    with st.chat_message(msg["role"]):
        if msg.get("filename"):
            st.caption(f"📎 {msg['filename']}")
        st.markdown(msg["content"])

if "selected_preset" not in st.session_state:
    st.session_state.selected_preset = None
if "compose_key_gen" not in st.session_state:
    st.session_state.compose_key_gen = 0
if "uploader_key_gen" not in st.session_state:
    st.session_state.uploader_key_gen = 0


def _on_text_change():
    st.session_state.selected_preset = None


presets = requests.get(f"{BACKEND_URL}/presets").json()
presets_by_id = {p["id"]: p for p in presets}

have_conversation = st.session_state.conversation_id is not None

if have_conversation:
    uploaded_file = None
else:
    uploaded_file = st.file_uploader(
        "Attach ledger CSV",
        type=["csv"],
        key=f"uploader_{st.session_state.uploader_key_gen}",
    )

can_compose = uploaded_file or have_conversation

if can_compose:
    if uploaded_file:
        st.caption("Choose an analysis, or type your own question below.")
    else:
        st.caption("Continue the conversation below.")

    st.caption("Ready-made questions: click to select")
    for i, preset in enumerate(presets, start=1):
        expanded_key = f"preset_expanded_{preset['id']}"
        if expanded_key not in st.session_state:
            st.session_state[expanded_key] = False

        select_col, toggle_col = st.columns([5, 1])
        with select_col:
            if st.button(
                f"{i}. {preset['label']}",
                key=f"preset_{preset['id']}",
                use_container_width=True,
            ):
                st.session_state.selected_preset = preset["id"]
        with toggle_col:
            toggle_label = "Hide" if st.session_state[expanded_key] else "Show"
            if st.button(toggle_label, key=f"toggle_{preset['id']}", use_container_width=True):
                st.session_state[expanded_key] = not st.session_state[expanded_key]

        if st.session_state[expanded_key]:
            st.text(preset["text"])

    if st.session_state.selected_preset:
        st.info(f"Selected: {presets_by_id[st.session_state.selected_preset]['label']}")

    compose_key = f"compose_text_{st.session_state.compose_key_gen}"
    st.text_area(
        "Or type your own question (typing here deselects any ready-made question above)",
        key=compose_key,
        on_change=_on_text_change,
    )

    send_clicked = st.button("Send", type="primary")
else:
    st.caption("Attach a ledger CSV to start.")
    send_clicked = False

if send_clicked:
    if st.session_state.selected_preset:
        query = presets_by_id[st.session_state.selected_preset]["text"]
    else:
        query = st.session_state.get(compose_key, "").strip()

    if not query:
        st.warning("Select an analysis or type a question first.")
    else:
        display_label = (
            presets_by_id[st.session_state.selected_preset]["label"]
            if st.session_state.selected_preset
            else query
        )

        st.session_state.messages.append(
            {
                "role": "user",
                "content": display_label,
                "filename": uploaded_file.name if uploaded_file else "",
            }
        )
        with st.chat_message("user"):
            if uploaded_file:
                st.caption(f"📎 {uploaded_file.name}")
            st.markdown(display_label)

        with st.chat_message("assistant"):
            placeholder = st.empty()
            reply = ""
            data = {"query": query}
            if st.session_state.conversation_id:
                data["conversation_id"] = st.session_state.conversation_id
            files = (
                {"file": (uploaded_file.name, uploaded_file.getvalue(), "text/csv")}
                if uploaded_file
                else None
            )

            spinner = st.spinner("Thinking...")
            spinner.__enter__()
            spinner_active = True

            with requests.post(
                f"{BACKEND_URL}/chat", data=data, files=files, headers=headers, stream=True
            ) as r:
                event = None
                data_lines = []
                done = False
                for line in r.iter_lines():
                    if line is None:
                        continue
                    if line == b"":
                        event_data = "\n".join(data_lines)
                        data_lines = []
                        if event == "conversation":
                            st.session_state.conversation_id = event_data.strip()
                        elif event == "message":
                            if spinner_active:
                                spinner.__exit__(None, None, None)
                                spinner_active = False
                            reply += event_data
                            placeholder.markdown(reply)
                        elif event == "done":
                            done = True
                        elif event == "error":
                            st.error(event_data)
                            done = True
                        if done:
                            break
                        continue
                    line = line.decode("utf-8")
                    if line.startswith("event:"):
                        event = line.split(":", 1)[1].strip()
                    elif line.startswith("data:"):
                        data_lines.append(line.split(":", 1)[1])

            if spinner_active:
                spinner.__exit__(None, None, None)

        if reply:
            st.session_state.messages.append({"role": "assistant", "content": reply})
        st.session_state.selected_preset = None
        st.session_state.compose_key_gen += 1
        if uploaded_file:
            st.session_state.uploader_key_gen += 1
        st.rerun()
