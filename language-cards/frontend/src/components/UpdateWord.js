import React, {useState} from 'react'
import updateWord from '../services/updateService'


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
            Update Word
        </button>
    </form>
    {message && <p>{message}</p>}
</div>

    );
  }
  
  export default UpdateWord;
  