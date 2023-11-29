Takes in path to the folder, list all the files along with number lines in them. 

Parameters:
```bash
-ext "comma separated extensions(default: '')" \
-dir "name of the directory(default: cwd)" \
-save "name of the file on which you want to save(default: not saving)" \
-disp=true/false "either to show it on stdout or not(default: true)"
```

Working features:
  * Can count only certain file types if their extensions are provided
  * Save to file option 
  * Don't handle directories unnecessarily
  * handle directories recursively

To do:
  * How to handle folder containing large number of files?
  * Handle the total better?
  * Make it run on windows as well