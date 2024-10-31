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


const Home = () => {
    const [words, setWords] = useState([]);
    const [sets, setSets] = useState([])
    const [currentUser, setCurrentUser] = useState(null)
    const [currentSet, setCurrentSet] = useState(null)
    const [showNext, setShowNext] = useState(false)
    const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
    const [showDefinition, setShowDefinition] = useState(true);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [isLoggedIn, setIsLoggedIn] = useState(false)
    const [isRegistering, setIsRegistering] = useState(false)
    const [unseenWords, setUnseenWords] = useState([])
    const token = localStorage.getItem('token');
    
    console.log("Home rendering with showDefinition:", showDefinition);

    const handleShowDefinitionClick = () => {
        setShowDefinition(prev => {
            const newShowDefinition = !prev;
            setShowNext(!prev); 
            return newShowDefinition;
        });
    };
    
    useEffect(() => {
        if (Array.isArray(words)) {
            setUnseenWords(words.filter(wordObj => !wordObj.seen));
        } else {
            setUnseenWords([]);
        }
    }, [words]);
    
    useEffect(() => {
        if (currentIndex >= words.length) {
            setCurrentIndex(Math.max(words.length - 1, 0));
        }
    }, [words, currentIndex]);

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
    

    useEffect(() => {
        if (isLoggedIn && currentSet && currentUser) {
            const fetchWords = async () => {
                try {
                    console.log("in fetch words currentSet:", currentSet);
                    console.log("in fetch words currentUser:", currentUser);
                    const fetchedWords = await showAllWords(currentUser.userID, currentSet);
                    console.log("Fetched words:", fetchedWords);
                    setWords(fetchedWords || []);
                } catch (error) {
                    console.error("Error fetching words:", error);
                    setWords([]); 
                }
            };
    
            fetchWords();
        }
    }, [isLoggedIn, currentSet, currentUser]);
    const handleLoginSuccess = () => {
        setIsLoggedIn(true)
    }

    const getFlagImagePath = (setName) => {
        const formattedName = setName.replace(/\s+/g, '').toLowerCase();
        return `/images/${formattedName}.png`; 
      };
    const handleLogout = () => {
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        setCurrentUser(null);
    };

    const handleCorrectNext = () => {
        const currentWord = words[currentIndex]?.word;
    
        if (unseenWords.length === 0) {
            return 0; 
        }

        if (currentWord) {
            updateSeen(currentWord, currentUser.userID, currentSet.setId)
                .then(() => {
                    handleSeen(currentWord);
                })
                .catch(error => {
                    console.error("Error updating seen status:", error);
                });
        }
        setShowDefinition(false);
        setShowNext(false);

        setCurrentIndex(prevIndex => {
            let nextIndex = (prevIndex + 1) % unseenWords.length;
            return words.findIndex(wordObj => wordObj.word === unseenWords[nextIndex].word);
        });
    };
    

    const handleIncorrectNext = () => {
    
        if (unseenWords.length === 0) {
            return 0; 
        }
        setShowDefinition(false);
        setShowNext(false);

        setCurrentIndex(prevIndex => {
            let nextIndex = (prevIndex + 1) % unseenWords.length;
            return words.findIndex(wordObj => wordObj.word === unseenWords[nextIndex].word);
        });
    };
    
    const handleSeen = (wordToMarkAsSeen) => {
        console.log("Words before update:", words); 
        const updatedWords = words.map(wordObj => {
            if (wordObj.word === wordToMarkAsSeen) {
                return { ...wordObj, seen: true };
            }
            return wordObj;
        });
        console.log("Words after update:", updatedWords);
        console.log(unseenWords)
        setWords(updatedWords);
    };
    

    const handleDelete = (wordToDelete) => {
        console.log("Deleting word:", wordToDelete);
        console.log("Current words state:", words);
    
        deleteWord(wordToDelete, currentUser.userID, currentSet.setId)
            .then(() => {
                const updatedWords = words.filter((word) => word.word !== wordToDelete);
                setWords(updatedWords);
                if (currentIndex >= updatedWords.length) {
                    setCurrentIndex(updatedWords.length - 1);
                }
            })
            .catch(error => {
                console.error('Error deleting word', error);
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

    const handleSetSelection = (set) => {
        setCurrentSet(set)
    }
   
     
    if (!isLoggedIn) {
        return (
            <div>
            <Header />
            <div className="max-w-[400px] mx-auto p-5 border border-gray-200 rounded-lg shadow-sm">
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
        <div className="flex flex-col items-center justify-center min-h-[40vh] mb-10">
            {currentUser ? (
                <div>
                    <Modal isOpen={isCreateModalOpen} onClose={() => setIsCreateModalOpen(false)}>
                        <CreateSets currentUser={currentUser} onClose={() => setIsCreateModalOpen(false)} />
                    </Modal>
                    {!currentSet ? (
                        <div>
                        <h1>Welcome {currentUser.username}, You have {sets.length} sets!</h1>
                        <div className="grid grid-cols-3 gap-2.5 p-2.5 items-center justify-center">                                {sets.map(set => (
                                 <button key={set.setId} className="border-2 border-black bg-none p-0 cursor-pointer block w-full h-full hover:brightness-110" onClick={() => handleSetSelection(set)}>
                                    <img src={getFlagImagePath(set.setName)} alt={`${set.setName} flag`} />
                                        <span className="text-3xl font-bold">{set.setName}</span>
                                    </button>
                        ))}
                            </div>
                        </div>
                    ) : (
                        <div className="flex flex-col items-center justify-center w-4/5 min-h-[40vh] mb-12">
                            <h1 className="whitespace-nowrap max-w-full">You have {unseenWords.length} {unseenWords.length === 1 ? "word" : "words"} remaining today </h1>
                            <div className="bg-white shadow-md p-20 mx-auto my-2.5 rounded-lg w-4/5 max-w-[600px]">
                                <Flashcard
                                    word={unseenWords[currentIndex]?.word || ""}
                                    words={words}
                                    definition={unseenWords[currentIndex]?.definition || ""}
                                    showDefinition={showDefinition}
                                    setShowDefinition={setShowDefinition}
                                    onDelete={handleDelete}
                                    showNext={showNext}
                                    setShowNext={setShowNext}
                                    onShowDefinitionClick={handleShowDefinitionClick}
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
                <h1>Loading...</h1>
             )}
        </div>
    </div>
)};

export default Home;