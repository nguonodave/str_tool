package main

import goReloaded "goReloaded/text_tool"

/*
--> after all modifications have been done, append words that are not instructions, i.e not "(up), (cap, 4) etc", to an empty array, "new_array_results
--> write the sentence to the created file using io.WriteString
--> defer output_file.Close() is for file handling
*/
func main() {
	goReloaded.FinalTool()
}
