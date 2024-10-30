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
                <button type="bg-blue-500 text-white px-4 py-2.5 border-none rounded cursor-pointer transition-colors duration-300 hover:opacity-90 mr-2.5">Delete</button>
            </form>
            {message && <p>{message}</p>}
        </div>
    );
}

export default DeleteWord;
  