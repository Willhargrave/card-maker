import React from 'react';
import updateWord from '../../../../services/updateService';
import WordForm from '../WordForm';

function UpdateWord({ onUpdate, word, definition, ...props }) {
    const handleSubmit = async (newWord, newDefinition, userID, setId, originalWord) => {
        await updateWord(originalWord, newWord, newDefinition, userID, setId);
        onUpdate(originalWord, newWord, newDefinition);
    };

    return (
        <WordForm 
            initialWord={word}
            initialDefinition={definition}
            onSubmit={handleSubmit}
            buttonText="Update"
            {...props}
        />
    );
}

export default UpdateWord;