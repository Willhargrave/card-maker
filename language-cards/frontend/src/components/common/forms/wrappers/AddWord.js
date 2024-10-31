import React from 'react';
import addWord from '../../../../services/apiService';
import WordForm from '../WordForm';

function AddWord({ onAdd, ...props }) {
    const handleSubmit = async (word, definition, userID, setId) => {
        await addWord(word, definition, userID, setId);
        onAdd(word, definition, userID, setId);
    };

    return (
        <WordForm 
            onSubmit={handleSubmit}
            buttonText="Add"
            {...props}
        />
    );
}

export default AddWord;
  