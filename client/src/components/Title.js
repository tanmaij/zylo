import React from 'react'

import './Style.css';

/* Title component */
const Title = ({ room }) => {
    return <div className={"chatApp__convTitle clearfix"}>
        <img src={room.avatar}/>
        <span>
            <b>{room.name}</b>
        </span>
        </div>;
  };

export default Title