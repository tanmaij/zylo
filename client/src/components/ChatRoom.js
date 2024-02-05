import React from 'react'
import { useState } from 'react'

import ChatBox from './ChatBox'
import { detectURL } from '../utils/URL'

import './Style.css';

/* ChatRoom component - composed of multiple ChatBoxes */
const ChatRoom = ({sendMessage, messages, typing, resetTyping, isTyping, owner, room, isSending}) => {
      return (
        <div className={"chatApp__room"}>
          <ChatBox
            key={0}
            room={room}
            owner={owner.name}
            ownerAvatar={owner.avatar}
            sendMessage={sendMessage}
            typing={typing}
            resetTyping={resetTyping}
            messages={messages}
            isTyping={isTyping}
            isSending={isSending}/>
        </div>
      );
  };

export default ChatRoom