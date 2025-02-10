
CREATE TABLE IF NOT EXISTS Titles (
    TitleID varchar [36] NOT NULL PRIMARY KEY,
    MangadexID varchar [36] NOT NULL,
    Name varchar [128] NOT NULL,
    Format VARCHAR [32] NOT NULL,
    Status VARCHAR [16] NOT NULL,
    Origin VARCHAR [2],
    Description VARCHAR[255]
    Year INTEGER NOT NULL,
    LastReadedChapter VARCHAR[8] NOT NULL,
    LastAPIChapter VARCAHR[8] NOT NULL,
    Tags VARCHAR [128] NOT NULL,
    AddedDate TEXT NOT NULL
);
