## package htf
### Funktionen
```
- func GetExecutable(url string, dir string) (bool, error) 
- func GetFile(url string, dir string, xFile string, perm uint32) (bool, error)
- func GetDownloadFilesInfo(url string) (*DownloadFilesInfo, error)
- func (fl *DownloadFilesInfo) GetFileInfo(FileName string) (*DownloadFileInfo, error) 
- func (f *DownloadFileInfo) Download(toFile string) error
```