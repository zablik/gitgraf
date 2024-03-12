package model

type Stats struct {
    LinesAdded    int `bson:"lines_added"`
    LinesDeleted  int `bson:"lines_deleted"`
    FilesAdded    int `bson:"files_added"`
    FilesDeleted  int `bson:"files_deleted"`
    FilesModified int `bson:"files_modified"`
}