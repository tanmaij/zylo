import React from 'react'

import './Style.css';

/* MessageItem component - composed of a message and the sender's avatar */
const MessageItem = ({ owner, sender, senderAvatar, message }) => {
    const messagePosition = owner === sender ? 'chatApp__convMessageItem--right' : 'chatApp__convMessageItem--left';
  
    return (
      <div className={"chatApp__convMessageItem " + messagePosition + " clearfix"}>
        {owner === sender ? null : <img src={senderAvatar} alt={sender} className="chatApp__convMessageAvatar" />}
        <div className="chatApp__convMessageValue" dangerouslySetInnerHTML={{ __html: message }}></div>
      </div>
    );
  };
  

export default MessageItem