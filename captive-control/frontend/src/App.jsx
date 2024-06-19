// import React, { useState } from "react";
// import { login } from "../wailsjs/go/main/App";

// function App() {
//   const [username, setUsername] = useState("");
//   const [password, setPassword] = useState("");
//   const [message, setMessage] = useState("");

//   const handleLogin = async () => {
//     try {
//       const response = await login(username, password);
//       setMessage(response ? "Login successful!" : "Login failed.");
//     } catch (err) {
//       setMessage("Error occurred: " + err);
//     }
//   };

//   return (
//     <div id="app">
//       <h1>Login Application</h1>
//       <div>
//         <label htmlFor="username">Username:</label>
//         <input
//           value={username}
//           onChange={(e) => setUsername(e.target.value)}
//           type="text"
//           id="username"
//         />
//       </div>
//       <div>
//         <label htmlFor="password">Password:</label>
//         <input
//           value={password}
//           onChange={(e) => setPassword(e.target.value)}
//           type="password"
//           id="password"
//         />
//       </div>
//       <button onClick={handleLogin}>Login</button>
//       <p>{message}</p>
//     </div>
//   );
// }

// export default App;


import React, { useState } from "react";
import { Login } from "../wailsjs/go/main/App";

function App() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");

  const handleLogin = async () => {
    try {
      const response = await Login(username, password);
      setMessage(response ? "Login successful!" : "Login failed.");
    } catch (err) {
      setMessage("Error occurred: " + err);
    }
  };

  return (
    <div id="app">
      <h1>Login Application</h1>
      <div>
        <label htmlFor="username">Username:</label>
        <input
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          type="text"
          id="username"
        />
      </div>
      <div>
        <label htmlFor="password">Password:</label>
        <input
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          type="password"
          id="password"
        />
      </div>
      <button onClick={handleLogin}>Login</button>
      <p>{message}</p>
    </div>
  );
}

export default App;
