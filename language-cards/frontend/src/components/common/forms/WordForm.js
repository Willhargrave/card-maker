import React, { useState } from 'react';

function WordForm({ 
    initialWord = '', 
    initialDefinition = '', 
    onSubmit, 
    onClose, 
    buttonText, 
    currentUser, 
    currentSet 
}) {
    const [word, setWord] = useState(initialWord);
    const [definition, setDefinition] = useState(initialDefinition);
    const [message, setMessage] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await onSubmit(word, definition, currentUser.userID, currentSet.setId, initialWord);
            setMessage(`${word} ${buttonText.toLowerCase()}ed successfully!`);
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
                    {buttonText} Word
                </button>
            </form>
            {message && <p>{message}</p>}
        </div>
    );
}

export default WordForm;
