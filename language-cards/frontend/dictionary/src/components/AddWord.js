import React, {useState} from 'react'
import addWord from '../services/apiService';
import "../styles/modal.css"
import "../styles/form.css"
function AddWord({ onAdd, onClose, currentUser, currentSet }) {
  const [word, setWord] = useState('');
  const [definition, setDefinition] = useState('');
  const [message, setMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log(currentSet)
    try {
      console.log("Sending data:", { 
        word, 
        definition, 
        userID: currentUser?.userID, 
        setId: currentSet?.setId 
      });
        await addWord(word, definition, currentUser.userID, currentSet);
        console.log("Adding word with:", { word, definition, userID: currentUser.userID, setId: currentSet });
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
        <form onSubmit={handleSubmit} className='form-container'>
          <input type="text" value={word} onChange={(e) => setWord(e.target.value)} placeholder="Word" className='input-field'required/>
          <input type="text" value={definition} onChange={(e) => setDefinition(e.target.value)} placeholder="Definition" className='input-field' required />
          <button type="submit" className='submit-button'>Add Word</button>
        </form>
        {message && <p>{message}</p>}
      </div>
    );
  }
  
  export default AddWord;
  