# Advanced use

The documentation component supports some advanced ~~exploits~~ features.

## Custom ordering of articles
As mentioned in [Writing Documentation](#/documentation/1-writing-doc), the directory-structure is sorted alphabetically.
But the structure is sorted by **filename**, and not by title. 
This means that you can sort your articles by prefixing the filenames with for example a number, but omitting the number for the titles.
This way [Advanced use](#/documentation/2-advanced-use) comes after [Writing Documentation](#/documentation/1-writing-doc) in the sidebar overview,
even though `A` comes before `W` in the alphabet. 
And you will not see the number in the titles, only in the URLs. This is because the titles in the sidebars are read from the first line in each of the files.

![number-prefixed filenames](https://har.ald.im/g/eD3AQIGB.png)
