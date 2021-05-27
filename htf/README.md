## package htf
### Funktionen
```
- func GetFile(url string, dir string, xFile string, locFile string) (bool, error)
- func GetExecutableFile(url string, dir string, xFile string) (bool, error) 
- func GetExecutable(url string, dir string) (bool, error) 
- func GetDownloadFilesInfo(url string) (*DownloadFilesInfo, error)
- func (fl *DownloadFilesInfo) GetFileInfo(FileName string) (*DownloadFileInfo, error) 
- func (f *DownloadFileInfo) Download(toFile string) error
```