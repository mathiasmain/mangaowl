//import { useState } from 'react'

import './App.css'
import HomeButton from './components/HomeButton.tsx'
import MangaList from "./components/MangaList.tsx"

function App() {
  // className='conteiner bg-blue-950 font-bold'
  return (
    <>
        <HomeButton></HomeButton>
      <div >
        <MangaList></MangaList>
      </div>
    </>
  )
}

export default App
