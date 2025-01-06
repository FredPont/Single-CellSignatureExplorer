# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.

# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# Written by Frederic PONT.
#(c) Frederic Pont 2018-2021

# this program was inspired from shiny examples
# in particular those ones :
# https://plot.ly/r/shinyapp-explore-diamonds/
# https://shiny.rstudio.com/gallery/plot-interaction-selecting-points.html
# many thanks to the authors for these examples !

# Aknowledgments
# we thank Aviv Regev and Asaf Madi (Broad Institute) for sharing R scripts to draw density plots.
# I thank Rüçhan EKREN who improved drop down lists with the picker Input widget

# start the program : Rscript Run.R

library(ggplot2)
library(Cairo)   # For nicer ggplot2 output when deployed on Linux

# data format : samples in columns and criteria in rows
# caution :  duplicate 'row.names' are not allowed
# caution : "-" forbiden in row or colnames

# read file in data directory
files <- list.files("data/")
file <- files[1]

samples <- read.table(paste("data/",file[1], sep = "", collapse = NULL),header = TRUE, row.names = 1, sep = "\t")
criteria <- colnames(samples)
samples <- as.data.frame(samples)

# coordinates to draw t-sne plot
criteriaX<-"tSNE_1"
criteriaY<-"tSNE_2"

# round function
arrondi <- function (x) ifelse(x>0, ifelse(x-trunc(x)>=0.5,trunc(x)+1,trunc(x)), ifelse(x-trunc(x)<=-0.5,trunc(x)-1,trunc(x)))

########################
# graphical interface  #
########################

########## objects ###########
ui <- fluidPage(
    headerPanel("Single-Cell Signature Viewer"),
    sidebarPanel(
        pickerInput('z', 'Signatures', choices = criteria, selected = criteria, options=pickerOptions(liveSearch=T)),       # signature list
        plotOutput("plot2", height = 300, width = 300 ),                                                # small density plot
        uiOutput("scaleBar"),                                                                           # the dynamic scale is rendered in the server section (see bellow) to take min and max calculated from the user input
        uiOutput("densityBar"),
        pickerInput('criteriaX', 'X', choices = criteria, selected = criteria, options=pickerOptions(liveSearch=T)),         # X coordinates of the plot
        pickerInput('criteriaY', 'Y', choices = criteria, selected = criteria, options=pickerOptions(liveSearch=T)),         # X coordinates of the plot
        pickerInput('sortOrder', 'Sort dots', choices = c("none", "A -> Z", "Z -> A" ), selected = criteria, options=pickerOptions(liveSearch=T)),         # X coordinates of the plot   
        # the dynamic density scale is rendered in the server  section (see bellow) to take min and max calculated from the user input
        #pickerInput('MapDotsize' , 'Map Dot size', choices = c(1:10), selected = 1, options=pickerOptions(liveSearch=T)),         #
    ),
    mainPanel(
        fluidRow(
            column(width = 4,
                plotOutput("plot1", height = 700, width = 1000) # the main t-sne plot
            )
        ),
        fluidRow(
            actionButton("SavePlot","Save PNG"),    # the save plot in PNG button
            actionButton("saveSVG","Save SVG"),    # the save plot in PNG button
            pickerInput('MapDotsize' , 'Map Dot size', choices = c(1:10), selected = 1, options=pickerOptions(liveSearch=T))         #
        )
    )
)


########## rendering ###########
# compute the max value of the color scale slider
calcMax <- function(input){
    value = max(input)
    min = min(input)
    max = 1.3 * value
    if (value <= 0) {
       max = 0.3 * (value - min) + value
    }

    # print(paste("value = ", value))
    # print(paste("min = ", min))
    # print(paste("max = ", max))
    return( c(value, min, max))
}

# sortDF sort the dataframe before ploting to bring dots in foreground/backgound
sortDF <-function(input, samples){
     # sort the dataframe to reveal some dots
    if (input$sortOrder == "A -> Z"){
        samples <- samples[order(samples[,input$z], decreasing = FALSE), ]
    } else if (input$sortOrder == "Z -> A") {
        samples <- samples[order(samples[,input$z], decreasing = TRUE), ]
    }
    return(samples)
}


server <- function(input, output) {
    # color slider
    output$scaleBar <- renderUI({  
        scaleval = calcMax(samples[,input$z])
        valuev = scaleval[1]
        minv = scaleval[2]
        maxv = scaleval[3]
        sliderInput("colorSelect",
                "Color scale on data score range:",
                min = minv,
                max = maxv,
                value = valuev,
                )
    })
    # density slider
    output$densityBar <- renderUI({
    scaleval = calcMax(samples[,input$z])
    
    valuev = scaleval[1]
    min_expression = scaleval[2]
    

    med_expression=((input$colorSelect-min_expression)/1.5)+min_expression
    if (med_expression > valuev){
        med_expression <- valuev # when slider color value is not updated before density slider, the med_expression of previous computation is used and med_expression can be > maxv
    } else if (med_expression < min_expression) {
        med_expression <- min_expression
    }
    # print(paste("med_expression = ",  med_expression))
    sliderInput("densitySelect",
                "Density scale on data score range:",
                min = min_expression, # 
                max = valuev, # density cannot be > to max(input$z)
                value = med_expression, # 
                ) # 
    })

    # density plot
    output$plot2 <- renderPlot({
        plot(density(samples[,input$z]),  col = "blue", main="", xlab="signature score", ylab="frequency")
    })
    # scatter plot
    output$plot1 <- renderPlot({

    # sort the dataframe to reveal some dots
    samples <- sortDF(input, samples)

    jet.colors <- colorRampPalette(c("#00007F", "blue", "#007FFF", "cyan", "#7FFF7F", "yellow", "#FF7F00", "red", "#7F0000"))

    mdotsize <- strtoi(input$MapDotsize)    # convert the dot size to integer
    
    max_expression=max(samples[,input$z])
    min_expression=min(samples[,input$z])
    med_expression=((input$colorSelect-min_expression)/1.5)+min_expression
    max_tsneX=arrondi(max(samples[,input$criteriaX]))+1
    min_tsneX=arrondi(min(samples[,input$criteriaX]))-1
    max_tsneY=arrondi(max(samples[,input$criteriaY]))+1
    min_tsneY=arrondi(min(samples[,input$criteriaY]))-1
    ggplot(samples, aes_string(x = input$criteriaX, y = input$criteriaY)) +
    stat_density_2d(data=samples[samples[,input$z]>input$densitySelect,], aes(fill=..level..,alpha=..level..), geom='polygon', colour='darkgrey', n=50) +
    geom_jitter(aes(x = samples[,input$criteriaX], y = samples[,input$criteriaY], color = samples[,input$z]), size=mdotsize, alpha=0.7) +
    labs(color = "color scale") + # color legend
    scale_colour_gradientn(colours = jet.colors(10), limits=c(min_expression, input$colorSelect)) +
    xlab(input$criteriaX) +
    theme_bw() +
    theme(panel.grid.major = element_blank(), panel.grid.minor = element_blank()) +  # remove grid
    scale_alpha_continuous(range=c(0.1,0.5))+
    scale_x_continuous(limits = c(min_tsneX, max_tsneX))+
    scale_y_continuous(limits = c(min_tsneY, max_tsneY))+
    theme(
    plot.title = element_text(size = 14),
    axis.text.x = element_text(size = 14, color='black'),
    axis.text.y = element_text(size = 14, color='black'),
    strip.text.x = element_text(size = 14, color='black',angle=90),
    strip.text.y = element_text(size = 14, color='black',angle=90),
    axis.title.x = element_text(size = 14, color='black'),
    axis.title.y = element_text(size = 14, angle=90))+
    ylab(input$criteriaY)
  })

    # save PNG image
  observeEvent(input$SavePlot,{
    # sort the dataframe to reveal some dots
    samples <- sortDF(input, samples)

    mdotsize <- strtoi(input$MapDotsize)    # convert the dot size to integer

    # progress bar
    withProgress(message = 'Saving plot', {
        
    inputPlot <- jet.colors <- colorRampPalette(c("#00007F", "blue", "#007FFF", "cyan", "#7FFF7F", "yellow", "#FF7F00", "red", "#7F0000"))

    max_expression=max(samples[,input$z])
    min_expression=min(samples[,input$z])
   
    max_tsneX=arrondi(max(samples[,input$criteriaX]))+1
    min_tsneX=arrondi(min(samples[,input$criteriaX]))-1
    max_tsneY=arrondi(max(samples[,input$criteriaY]))+1
    min_tsneY=arrondi(min(samples[,input$criteriaY]))-1

    p<-ggplot(samples , aes_string(x = input$criteriaX, y = input$criteriaY)) +
    stat_density_2d(data=samples[samples[,input$z]>input$densitySelect,], aes(fill=..level..,alpha=..level..),geom='polygon',colour='darkgrey',n=50) +
    geom_jitter(aes(x = samples[,input$criteriaX], y = samples[,input$criteriaY], color=samples[,input$z]), size=mdotsize, alpha=0.7) +
    labs(color = "color scale") + # color legend
    scale_colour_gradientn(colours = jet.colors(10), limits=c(min_expression, input$colorSelect)) +
    xlab(input$criteriaX) +
    theme_bw() +
    theme(panel.grid.major = element_blank(), panel.grid.minor = element_blank()) +  # remove grid
    scale_alpha_continuous(range=c(0.1,0.5))+
    scale_x_continuous(limits = c(min_tsneX, max_tsneX))+
    scale_y_continuous(limits = c(min_tsneY, max_tsneY))+
    theme(
    plot.title = element_text(size = 14),
    axis.text.x = element_text(size = 14, color='black'),
    axis.text.y = element_text(size = 14, color='black'),
    strip.text.x = element_text(size = 14, color='black',angle=90),
    strip.text.y = element_text(size = 14, color='black',angle=90),
    axis.title.x = element_text(size = 14, color='black'),
    axis.title.y = element_text(size = 14,angle=90))+
    ylab(input$criteriaY)+
    labs(title = input$z)
    file<-paste0("plots/", input$z,".png")
    ggsave(file,p, width = 243, height = 150, units = "mm", dpi = 200)
    print("image saved in plots directory !")
    })
  })

    # save SVG image
  observeEvent(input$saveSVG,{
    # sort the dataframe to reveal some dots
    samples <- sortDF(input, samples)

    mdotsize <- strtoi(input$MapDotsize)    # convert the dot size to integer
    
    # progress bar
    withProgress(message = 'Saving plot', {

    inputPlot <- jet.colors <- colorRampPalette(c("#00007F", "blue", "#007FFF", "cyan", "#7FFF7F", "yellow", "#FF7F00", "red", "#7F0000"))

    
    max_expression=max(samples[,input$z])
    min_expression=min(samples[,input$z])
    
    max_tsneX=arrondi(max(samples[,input$criteriaX]))+1
    min_tsneX=arrondi(min(samples[,input$criteriaX]))-1
    max_tsneY=arrondi(max(samples[,input$criteriaY]))+1
    min_tsneY=arrondi(min(samples[,input$criteriaY]))-1
    p<-ggplot(samples , aes_string(x = input$criteriaX, y = input$criteriaY)) +
    stat_density_2d(data=samples[samples[,input$z]>input$densitySelect,], aes(fill=..level..,alpha=..level..),geom='polygon',colour='darkgrey',n=50) +
    geom_jitter(aes(x = samples[,input$criteriaX], y = samples[,input$criteriaY], color=samples[,input$z]), size=mdotsize, alpha=0.7) +
    labs(color = "color scale") + # color legend
    scale_colour_gradientn(colours = jet.colors(10), limits=c(min_expression, input$colorSelect)) +
    xlab(input$criteriaX) +
    theme_bw() +
    theme(panel.grid.major = element_blank(), panel.grid.minor = element_blank()) +  # remove grid
    scale_alpha_continuous(range=c(0.1,0.5))+
    scale_x_continuous(limits = c(min_tsneX, max_tsneX))+
    scale_y_continuous(limits = c(min_tsneY, max_tsneY))+
    theme(
    plot.title = element_text(size = 14),
    axis.text.x = element_text(size = 14, color='black'),
    axis.text.y = element_text(size = 14, color='black'),
    strip.text.x = element_text(size = 14, color='black',angle=90),
    strip.text.y = element_text(size = 14, color='black',angle=90),
    axis.title.x = element_text(size = 14, color='black'),
    axis.title.y = element_text(size = 14,angle=90))+
    ylab(input$criteriaY)+
    labs(title = input$z)
    file<-paste0("plots/", input$z,".svg")
    ggsave(file,p, width = 300, height = 200, units = "mm")
    print("image saved in plots directory !")
    })
  })

  # format date
    now <- function(){
        d <- format(Sys.time(), "%Y_%m_%d_%X")
        gsub(":","",d)
    }

  # save data in file
    save_data_flatfile <- function(data) {
        file_name <- paste("results/" , now(), ".csv")
        write.table(x = data, file =  file_name, sep = "\t",
        row.names = TRUE, quote = TRUE, col.names = NA)
    }


}

shinyApp(ui, server)
