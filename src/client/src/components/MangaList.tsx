
//import { useState, useEffect, useRef } from 'react';
import { useQuery } from '@tanstack/react-query';

interface Title {
  ID:                string,
  MangadexID:        string,
  Name:              string,
  Format:            string,
  Status:            string,
  Country:           string,
  Year:              number,
  LastReadedChapter: string,
  LastAPIChapter:    string,
  Tags:              string,
  AddedDate:         string,
}

function MangaList() { // Confuso? vá para: https://www.youtube.com/watch?v=00lxm_doFYw

  const {data, error, isError, isLoading} = useQuery({
    queryKey: ["titles"],
    queryFn: async () => {
      return (await (await fetch("http://127.0.0.1:5582/api/titles")).json()) as Title[];
    },
  });

  
  if (isError) {
    {console.log(error);}
    return <div>Something went wrong, please try again.</div>
  }
  return (
    <>
      <h1 className="font-bold text-3xl flex-auto border-4 mb-3">Get them:</h1>
      <div>
        
        { isLoading && <div>Loading...</div>}
        
        <ul className="">
            {data?.map((title) => {
                return <li key={title.AddedDate}>{"\nName: "}{title.Name}{"\nUser Last Chapter :"}{title.LastReadedChapter}{"\nMangaDex Last Chapter :"},{title.LastAPIChapter}{"\nTitle's tags :"}{title.Tags} </li>;  // Supondo que `title` tenha uma propriedade `name`
            })}
            
        </ul>
      </div>
    </>
    
  );
}

export default MangaList;