import React, { useState, useEffect } from 'react';

import './Style.css';

/* render the chatroom */
setTimeout(() => {
  ReactDOM.render(<ChatRoom />, document.getElementById('chatApp'));
}, 400);
