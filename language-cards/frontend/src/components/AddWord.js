import React, {useState} from 'react'
import addWord from '../services/apiService';

function AddWord({ onAdd, onClose, currentUser, currentSet }) {
  const [word, setWord] = useState('');
  const [definition, setDefinition] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log("currentSet:", currentSet);  
    console.log("setId being passed:", currentSet.setId); 

    try {
        await addWord(
          word, 
          definition, 
          currentUser.userID, 
          currentSet.setId
        );
        onAdd(word, definition, currentUser.userID, currentSet.setId);
        setMessage(`${word} added successfully!`);
        onClose();
    } catch (error) {
        console.error(error);
        setMessage(error.message);
    }
};

    return (
      <div>
      <form onSubmit={handleSubmit} className="bg-white p-5 rounded">
          <input 
              type="text" 
              value={word} 
              onChange={(e) => setWord(e.target.value)} 
              placeholder="Word" 
              className="w-[calc(100%-30px)] p-2.5 mb-4 border border-gray-300 rounded box-border"
              required
          />
          <input 
              type="text" 
              value={definition} 
              onChange={(e) => setDefinition(e.target.value)} 
              placeholder="Definition" 
              className="w-[calc(100%-30px)] p-2.5 mb-4 border border-gray-300 rounded box-border"
              required
          />
          <button type="submit" className="bg-blue-500 text-white px-4 py-2.5 border-none rounded cursor-pointer transition-colors duration-300 hover:opacity-90 mr-2.5">
              Add Word
          </button>
      </form>
      {message && <p>{message}</p>}
  </div>
    );
  }
  
  export default AddWord;
  