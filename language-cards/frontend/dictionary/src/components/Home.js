import React, { useEffect, useState } from "react";
import Flashcard from "./Flashcard";
import showAllWords from "../services/showServive";
import fetchUserSets from "../services/fetchSetService";
import updateSeen from "../services/seenService";
import Login from "./Login";
import Header from "./Header";
import deleteWord from "../services/deleteService";
import Register from "./Register";
import CreateSets from "./CreateSets";
import Modal from "./Modal";
import {jwtDecode} from 'jwt-decode';
import "../styles/home.css"
import "../styles/buttons.css"
const Home = () => {
    const [words, setWords] = useState([]);
    const [sets, setSets] = useState([])
    const [currentUser, setCurrentUser] = useState(null)
    const [currentSet, setCurrentSet] = useState(null)
    const [showNext, setShowNext] = useState(false)
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
    const [showDefinition, setShowDefinition] = useState(false);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [isRegistering, setIsRegistering] = useState(false)
    const unseenWords = Array.isArray(words) ? words.filter(wordObj => !wordObj.seen) : [];
    const token = localStorage.getItem('token');
    
    //current-index
    useEffect(() => {
        // Reset currentIndex if it becomes out of bounds
        if (currentIndex >= words.length) {
            setCurrentIndex(Math.max(words.length - 1, 0));
        }
    }, [words, currentIndex]);
    //current-user
    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            try {
                const decodedToken = jwtDecode(token); 
                setIsLoggedIn(true);
                setCurrentUser({ 
                    username: decodedToken.username,
                    userID: decodedToken.UserID 
                });
                console.log("Current User token:", decodedToken, decodedToken.username, decodedToken.UserID); 
                console.log("Current User:", currentUser)
            } catch (error) {
                console.error('Error decoding the token:', error);
            }
        }
    }, [token]);
    //fetch-sets
    useEffect(() => {
        if (currentUser && currentUser.userID) {
            console.log("User set in fetchSets", currentUser.userID)
            const fetchSets = async () => {
                console.log("Fetching sets for user:", currentUser.UserID);
                try {
                    const userSets = await fetchUserSets(currentUser.UserID);
                    console.log("Fetched sets:", userSets);
                    setSets(userSets || []);
                } catch (error) {
                    console.error("Error fetching sets:", error);
                }
            };
            fetchSets();
        }
    }, [currentUser]);
    

    useEffect(() => {
        console.log("Current User updated:", currentUser);
    }, [currentUser]);
    
    //fetch-words
    useEffect(() => {
        if (isLoggedIn && currentSet && currentUser) {
            const fetchWords = async () => {
                try {
                    console.log("in fetch words currentSet:", currentSet);
                    console.log("in fetch words currentUser:", currentUser);
                    const fetchedWords = await showAllWords(currentUser.userID, currentSet);
                    console.log("Fetched words:", fetchedWords);
                    setWords(fetchedWords || []); // Set to empty array if null
                } catch (error) {
                    console.error("Error fetching words:", error);
                    setWords([]); // Set to empty array in case of error
                }
            };
    
            fetchWords();
        }
    }, [isLoggedIn, currentSet, currentUser]);
    const handleLoginSuccess = () => {
        setIsLoggedIn(true)
    }
    //flag
    const getFlagImagePath = (setName) => {
        // Replace spaces and lowercase
        const formattedName = setName.replace(/\s+/g, '').toLowerCase();
        return `/images/${formattedName}.png`; // Adjust path based on your folder structure
      };
    const handleLogout = () => {
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        setCurrentUser(null);
    };
    const handleCorrectNext = () => {
        const currentWord = words[currentIndex]?.word;
    
        if (unseenWords.length === 0) {
            return 0; // No unseen words left, reset to 0
        }
    
        // Update the seen status of the current word
        if (currentWord) {
            updateSeen(currentWord, currentUser.userID, currentSet)
                .then(() => {
                    handleSeen(currentWord);
                })
                .catch(error => {
                    console.error("Error updating seen status:", error);
                });
        }
        setShowDefinition(false);
        setShowNext(false);
    
        // Move to the next word
        setCurrentIndex(prevIndex => {
            let nextIndex = (prevIndex + 1) % unseenWords.length;
            return words.findIndex(wordObj => wordObj.word === unseenWords[nextIndex].word);
        });
    };
    

    const handleIncorrectNext = () => {
    
        if (unseenWords.length === 0) {
            return 0; // No unseen words left, reset to 0
        }
        setShowDefinition(false);
        setShowNext(false);

        setCurrentIndex(prevIndex => {
            let nextIndex = (prevIndex + 1) % unseenWords.length;
            return words.findIndex(wordObj => wordObj.word === unseenWords[nextIndex].word);
        });
    };
    
    const handleSeen = (wordToMarkAsSeen) => {
        console.log("Words before update:", words); // Before update
        const updatedWords = words.map(wordObj => {
            if (wordObj.word === wordToMarkAsSeen) {
                return { ...wordObj, seen: true };
            }
            return wordObj;
        });
        console.log("Words after update:", updatedWords); // After update
        console.log(unseenWords)
        setWords(updatedWords);
    };
    

    const handleDelete = (wordToDelete) => {
        console.log("Deleting word:", wordToDelete);
        console.log("Current words state:", words);
    
        deleteWord(wordToDelete, currentUser.userID, currentSet)
            .then(() => {
                const updatedWords = words.filter((word) => word.word !== wordToDelete);
                setWords(updatedWords);
                if (currentIndex >= updatedWords.length) {
                    setCurrentIndex(updatedWords.length - 1);
                }
            })
            .catch(error => {
                console.error('Error deleting word', error);
                // Handle delete error
            });
    };
    

    const handleUpdate = (originalWord, updatedWord, updatedDefinition) => {
        const updatedWords = words.map((wordObj) => {
            if (wordObj.word === originalWord) {
                return { ...wordObj, word: updatedWord, definition: updatedDefinition};
            }
            return wordObj;
        });
        setWords(updatedWords);
    };

    const handleAddWord = (newWord, newDefinition) => {
        const newWordsArray = [...words, { word: newWord, definition: newDefinition }];
        setWords(newWordsArray);
    };

    const handleSetSelection = (setId) => {
        setCurrentSet(setId)
    }
   
     
    if (!isLoggedIn) {
        return (
            <div>
            <Header />
            <div className="login-container">
                {isRegistering ? (
                    <Register onRegister={() => setIsRegistering(false)} setIsLoggedIn={setIsLoggedIn} onLoginSuccess={handleLoginSuccess} />
                ) : (
                    <>
                        <Login onLogin={() => setIsLoggedIn(true)} setIsLoggedIn={setIsLoggedIn} loginSuccess={handleLoginSuccess}/>
                        <button onClick={() => setIsRegistering(true)} className="register-button">Register</button>
                    </>
                )}
            </div>
            </div>
        );
        
    }

    return (
        <div>
        <Header handleLogout={handleLogout} isLoggedIn={isLoggedIn} setIsCreateModalOpen={setIsCreateModalOpen}/>
        <div className="home-container">
            {currentUser ? (
                <div>
                    <Modal isOpen={isCreateModalOpen} onClose={() => setIsCreateModalOpen(false)}>
                        <CreateSets currentUser={currentUser} onClose={() => setIsCreateModalOpen(false)} />
                    </Modal>
                    {!currentSet ? (
                        <div>
                        <h1>Welcome {currentUser.username}, You have {sets.length} sets!</h1>
                            <div className="set-container">
                                {sets.map(set => (
                                 <button key={set.setId} className="image-button" onClick={handleSetSelection}>
                                    <img src={getFlagImagePath(set.setName)} alt={`${set.setName} flag`} />
                                        <span className="card-title">{set.setName}</span>
                                    </button>
                        ))}
                            </div>
                        </div>
                    ) : (
                        <div className="main-container">
                                <h1>You have {unseenWords.length} {unseenWords.length === 1 ? "word" : "words"} remaining today </h1>
                             <div className="card">
                                <Flashcard
                                    word={words[currentIndex]?.word || ""}
                                    words={words}
                                    definition={words[currentIndex]?.definition || ""}
                                    showDefinition={showDefinition}
                                    setShowDefinition={setShowDefinition}
                                    onDelete={handleDelete}
                                    showNext={showNext}
                                    setShowNext={setShowNext}
                                    onUpdate={handleUpdate}
                                    onAdd={handleAddWord}
                                    onCorrect={handleCorrectNext}
                                    onIncorrect={handleIncorrectNext}
                                    unseenWords={unseenWords}
                                    sets={sets}
                                    setCurrentSet={setCurrentSet}
                                    currentUser={currentUser}
                                    currentSet={currentSet}
                                />
                            </div>
                        </div>
                    )}
                </div>
            ) : (
                <h1>Loading...</h1> // or any placeholder you prefer
            )}
        </div>
        </div>
)};

export default Home;