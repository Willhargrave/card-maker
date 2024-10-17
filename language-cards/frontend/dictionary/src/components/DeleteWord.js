import React, {useState} from 'react'
import deleteWord from '../services/deleteService'

function DeleteWord({ currentUser, currentSet }) {
    const [word, setWord] = useState('')
    const [message, setMessage] = useState('')

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await deleteWord(word, currentUser.userID, currentSet.setId);
            setMessage(`${word} deleted successfully!`);
            setWord('')
        } catch (error) {
            console.error('Error deleting word:', error);
            setMessage(`Cannot delete ${word} as it is not in your deck`);
        }
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <input type="text" value={word} onChange={(e) => setWord(e.target.value)} placeholder="Word" required />
                <button type="submit">Delete</button>
            </form>
            {message && <p>{message}</p>}
        </div>
    );
}

export default DeleteWord;
  