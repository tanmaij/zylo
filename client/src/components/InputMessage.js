import React from 'react'
import { useState } from 'react'

import './Style.css';

import 'material-icons/iconfont/material-icons.css';

/* InputMessage component - used to type the message */
const InputMessage = ({ isSending, owner, ownerAvatar, sendMessage, typing, resetTyping }) => {
    const [messageInput, setMessageInput] = useState('');
  
    const handleSendMessage = (event) => {
      event.preventDefault();
      if (messageInput.length > 0) {
        sendMessage(owner, ownerAvatar, messageInput);
        setMessageInput('');
      }
    };
  
    const handleTyping = (event) => {
    };

    const loadingClass = isSending ? 'chatApp__convButton--loading' : '';
    const sendButtonIcon = <i className={"material-icons"}>send</i>;
  
    return (
      <form onSubmit={handleSendMessage}>
        <input type="hidden" value={owner} />
        <input type="hidden" value={ownerAvatar} />
        <input
          type="text"
          value={messageInput}
          onChange={(e) => setMessageInput(e.target.value)}
          className={"chatApp__convInput"}
          placeholder="Text message"
          onKeyDown={handleTyping}
          onKeyUp={handleTyping}
          tabIndex="0"
        />
        <div className={'chatApp__convButton ' + loadingClass} onClick={handleSendMessage}>
          {sendButtonIcon}
        </div>
      </form>
    );
  };
  

export default InputMessage