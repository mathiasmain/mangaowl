
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

export function COuntry(country: string): string{
  if (country === "ko") {
    return "South Korea"
  }
  else if (country === "ja") {
    return "Japan"
  } else if (country === "ch")
    {
      return "China"
    }
  else {
    return "Country not reconized."
  }
}


function MangaList() { // Confuso? vá para: https://www.youtube.com/watch?v=00lxm_doFYw

  const {data, error, isError, isLoading} = useQuery({
    queryKey: ["titles"], // text-center flex-col flex-auto
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
      <div>
        
        { isLoading && <div>Loading...</div>}
        
        <ul className=''>
            {data?.map((title) => {
                if (title.ID === null) {
                  return <p>Your list is empty! Search and save titles to start your marker list.</p>
                }
                return <li className = "bg-gray-700 box-border border-2 rounded-2xl m-2 max-w-auto" key={title.AddedDate}>
                  <div className="font-bold text-yellow-300">{title.Name}</div>
                  <div className="flex space-x-2 mr-8 ml-8 ">
                    <div className="pl-2 pr-2">
                        <div className="flex space-x-2 ml-8 mr-8 mb-2 mt-2">
                            <div className="pr-2 pl-2 text-white rounded-lg">{"User's Chapter : "+title.LastReadedChapter}</div>
                            <div className="pr-2 pl-2 text-white rounded-lg">{"Latest Chapter: "+ title.LastAPIChapter}</div>
                        </div>
                        <div className="flex space-x-2 ml-8 mr-8 mb-2 mt-2">{title.Tags.split(",").map((tag) => {
                          return <div className="pr-2 pl-2 text-white rounded-lg">{tag}</div>
                        })}</div>
                        <div className="flex space-x-2  ml-8 mr-8 mb-2 mt-2">
                          <div className="pr-2 pl-2 text-white rounded-lg">{"Release year: "+title.Year}</div>
                          <div className="pr-2 pl-2 text-white rounded-lg">{"Original Country: "+COuntry(title.Country)}</div>
                          <div className="pr-2 pl-2 text-white rounded-lg">{"Format: "+title.Format}</div>
                          <div className="pr-2 pl-2 text-white rounded-lg">{"Status: "+title.Status}</div>
                        </div>
                    </div>
                  </div>
                  </li>;
            })}
            
        </ul>
      </div>
    </>
    
  );
}

export default MangaList;