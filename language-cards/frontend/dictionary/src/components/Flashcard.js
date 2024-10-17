import React, {useState} from "react";
import UpdateWord from "./UpdateWord";
import AddWord from "./AddWord";
import Modal from "./Modal";
import "../styles/buttons.css"
import "../styles/flashcard.css"
function Flashcard({word, definition, showDefinition, setShowDefinition, onDelete, onCorrect, onIncorrect, showNext, setShowNext, onUpdate, onAdd, unseenWords, currentUser, currentSet, setCurrentSet}) {

    const [isAddModalOpen, setIsAddModalOpen] = useState(false)
    const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false)

    const handleShowDefinitionClick = () => {
        setShowDefinition(!showDefinition);
        setShowNext(!showNext);
    };

    const handleDeleteClick = () => {
        onDelete(word)
    }

    return (
        <div>
            {unseenWords.length === 0 ? (
                 <>
                 <p>Well done! You've finished for today!</p>
                 <button className="button add-button" onClick={() => setIsAddModalOpen(true)}>Add New Word</button>
                 <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
                 <AddWord onAdd={onAdd} currentUser={currentUser} currentSet={currentSet} onClose={() => setIsAddModalOpen(false)}/>
                 </Modal>
             </>
            ) : (
                <>
            <div className="buttons-container">
            <button className="button add-button" onClick={() => setIsAddModalOpen(true)}>Add</button>
            <button className="button update-button" onClick={() => setIsUpdateModalOpen(true)}>Update</button>
            <button className="button delete-button" onClick={handleDeleteClick}>Delete</button>
            <button className="button return-button"onClick={() => setCurrentSet(false)}>Return</button>
            </div>
            <Modal isOpen={isAddModalOpen} onClose={() => setIsAddModalOpen(false)}>
                <AddWord onAdd={onAdd} currentUser={currentUser} currentSet={currentSet} onClose={() => setIsAddModalOpen(false)}/>
            </Modal>
            <Modal isOpen={isUpdateModalOpen} onClose={() => setIsUpdateModalOpen(false)}>
                <UpdateWord word={word} definition={definition} onUpdate={onUpdate} currentSet={currentSet} currentUser={currentUser} onClose={() => setIsAddModalOpen(false)} />
            </Modal>

            <div className="word-container">
            <h2 className="displayed-word">{word}</h2>
            {showDefinition && <p>{definition}</p>}
            <button className="button show-definition-button"onClick={handleShowDefinitionClick}>
               {showDefinition ? 'Hide Definition' : 'Show Definition'}
            </button>
            {showNext && (
                <>
                    <button className="button correct-button" onClick={() => onCorrect(word)}>Correct</button>
                    <button className="button wrong-button" onClick={onIncorrect}>Wrong</button>
                </>
            )}
            </div>
            </>
    )}
        </div>
    );
}

export default Flashcard;