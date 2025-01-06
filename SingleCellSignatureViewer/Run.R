# This function install R packages that are missing to run the software
# credits : https://stackoverflow.com/questions/9341635/check-for-installed-packages-before-running-install-packages
using<-function(...) {
    libs<-unlist(list(...))
    req<-unlist(lapply(libs,require,character.only=TRUE))
    need<-libs[req==FALSE]
    if(length(need)>0){ 
        install.packages(need)
        lapply(need,require,character.only=TRUE)
    }
}

using("shiny","ggplot2","Cairo","svglite","shinyWidgets")

# start interactive plot
library(shiny)
library(shinyWidgets)
print("caution : sign -, +, (), ', :, and special characters as well are forbiden in row or colnames !")
print("rownames or colnames should not start by a number !")
runApp("SingleCellSignatureVewer-055.R")
