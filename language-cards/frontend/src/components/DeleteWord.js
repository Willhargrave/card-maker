import React, {useState} from 'react'
import deleteWord from '../services/deleteService'

function DeleteWord({ currentUser, currentSet }) {
    const [word, setWord] = useState('')
    const [message, setMessage] = useState('')
    
    const handleSubmit = async (e) => {
        e.preventDefault();
        
        // Add validation for empty word
        if (!word.trim()) {
            setMessage('Please enter a word to delete');
            return;
        }

        try {
            await deleteWord(word.trim(), currentUser.userID, currentSet.setId);
            setMessage(`${word} deleted successfully!`);
            setWord('');
        } catch (error) {
            console.error('Error deleting word:', error);
            setMessage(`Cannot delete ${word} as it is not in your deck`);
        }
    };

    return (
        <div className="p-4">
            <form onSubmit={handleSubmit} className="flex gap-4">
                <input 
                    type="text" 
                    value={word} 
                    onChange={(e) => setWord(e.target.value)} 
                    placeholder="Word" 
                    className="w-full p-2.5 border border-gray-300 rounded"
                    required 
                />
                <button 
                    type="submit" 
                    className="bg-blue-500 text-white px-4 py-2.5 border-none rounded cursor-pointer transition-colors duration-300 hover:opacity-90"
                >
                    Delete
                </button>
            </form>
            {message && <p className="mt-4 text-center">{message}</p>}
        </div>
    );
}

export default DeleteWord;