// https://veryfirefly.github.io/xfs-filesystem-problem/

#include <stdio.h>
#include <dirent.h>
#include <sys/stat.h>

int main(int argc, char *argv[])
{
  struct stat info;
  DIR *dirp;
  struct dirent* dent;

  //If no args
  if (argc == 1)
  {

    argv[1] = ".";
    dirp = opendir(argv[1]); // specify directory here: "." is the "current directory"
    do
    {
      dent = readdir(dirp);
      if (dent)
      {
        printf("%c ", dent->d_type);
        printf("%s \n", dent->d_name);

        /* if (!stat(dent->d_name, &info))
         {
         //printf("%u bytes\n", (unsigned int)info.st_size);

         }*/
      }
    } while (dent);
    closedir(dirp);
  }

  //If specified directory 
  if (argc > 1)
  {
    dirp = opendir(argv[1]); // specify directory here: "." is the "current directory"
    do
    {
      dent = readdir(dirp);
      if (dent)
      {
        switch (dent -> d_type)
        {
          case DT_DIR:
            printf("DIR: %d, %s\n", dent -> d_type, dent->d_name);
            break;
          case DT_REG:
            printf("FILE: %d, %s\n", dent -> d_type, dent->d_name);
            break;
          default:
            printf("Other: %d, %s\n", dent -> d_type, dent->d_name);
            break;
        }

        // printf("%c ", dent->d_type);
        // printf("%s \n", dent->d_name);
        /*  if (!stat(dent->d_name, &info))
         {
         printf("%u bytes\n", (unsigned int)info.st_size);
         }*/
      }
    } while (dent);
    
    closedir(dirp);
  }
  return 0;
}
