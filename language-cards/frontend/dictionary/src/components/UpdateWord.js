import React, {useState} from 'react'
import updateWord from '../services/updateService'
import "../styles/form.css"
import "../styles/modal.css"

function UpdateWord({word: initialWord, definition: initialDefinition, onClose, onUpdate, currentUser, currentSet}) {
   
    const [word, setWord] = useState(initialWord)
    const [definition, setDefinition] = useState(initialDefinition)
    const [message, setMessage] = useState('')

    const handleSubmit = async (e) => {
      e.preventDefault();
      try {
          await updateWord(initialWord, word, definition, currentUser.userID, currentSet.setId);
          setMessage(`${word} updated successfully! New definition: ${definition}`);
          onUpdate(initialWord, word, definition); 
          onClose(); 
      } catch (error) {
          setMessage(error.message); 
      }
  };

    return (
        <div>
        <form onSubmit={handleSubmit} className='form-container'>
          <input type="text" value={word} onChange={(e) => setWord(e.target.value)} placeholder="Word" className='input-field'required/>
          <input type="text" value={definition} onChange={(e) => setDefinition(e.target.value)} placeholder="Definition" className='input-field' required/>
          <button type="submit" className='submit-button'>Update Word</button>
        </form>
        {message && <p>{message}</p>}
      </div>
    );
  }
  
  export default UpdateWord;
  