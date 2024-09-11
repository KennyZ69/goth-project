const chatId = document.getElementById("chat_id")
const ws = new WebSocket(`ws://http://127.0.0.1:7331/ws?chat_id=${chatId}`)

// When a message is received from the WebSocket server
ws.onmessage = function(event) {
  const message = JSON.parse(event.data);
  appendMessageToChat(message);
};

// Function to append a message to the chat
function appendMessageToChat(message) {
  const chatRoom = document.getElementById('chat_room');

  // Create a new list item for the message
  const messageElement = document.createElement('li');
  messageElement.classList.add('p-3', 'rounded-lg', 'shadow', 'max-w-xs');

  if (message.Client_id == parseInt(`{ data.CurrentUser.Id }`)) {
    messageElement.classList.add('bg-blue-100', 'text-right');
  } else {
    messageElement.classList.add('bg-gray-100', 'text-left');
  }

  messageElement.innerHTML = `
        <div class="text-sm text-gray-800">${message.text}</div>
        <div class="text-xs text-gray-500 mt-1">${new Date(message.created_at).toLocaleString()}</div>
    `;

  chatRoom.appendChild(messageElement);

  // Scroll to the bottom after appending
  chatRoom.scrollTop = chatRoom.scrollHeight;
}

// Handle sending the message via WebSocket
document.getElementById('chatForm').addEventListener('submit', function(event) {
  event.preventDefault();

  const messageInput = document.getElementById('chatInput');
  const messageText = messageInput.value;

  if (messageText.trim() != "") {
    const message = {
      chat_id: chatId,
      text: messageText
    };

    // Send the message via WebSocket
    ws.send(JSON.stringify(message));

    // Clear the input
    messageInput.value = '';
  }
});
