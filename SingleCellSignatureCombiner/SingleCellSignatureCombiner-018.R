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
#(c) Frederic Pont 2018

# this program was inspired from shiny examples
# in particular those ones :
# https://plot.ly/r/shinyapp-explore-diamonds/
# https://shiny.rstudio.com/gallery/plot-interaction-selecting-points.html
# many thanks to the authors for these examples !

# Aknowledgments
# we thank Aviv Regev and Asaf Madi (Broad Institute) for sharing R scripts to draw density plots.
# I thank Rüçhan EKREN who improved drop down lists with the picker Input widget

# start the program : Rscript Run.R

# function to normalise vector x between 0-1
norm <- function(x){
    a=(x-min(x))/(max(x)-min(x))
   return(a)
}



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
    headerPanel("Single-Cell Signature Combiner"),
    sidebarPanel(
        pickerInput('z', 'Signatures 1', choices = criteria, selected = criteria, options=pickerOptions(liveSearch=T)),              # signature list 1
        selectInput('w', 'Signatures 2', choices = criteria, selectize=FALSE, selected = criteria),              # signature list 2
        selectInput('o', 'Operator', choices = c("+", "-", "*"), selectize=FALSE, selected = c("-")),            # signature list 2
        plotOutput("plot2", height = 300, width = 300 ),                                                         # small density plot
        uiOutput("scaleBar"),                                                                                    # the dynamic scale is rendered in the server section (see bellow) to take min and max calculated from the user input
        uiOutput("densityBar"),                                                                                  # the dynamic density scale is rendered in the server  section (see bellow) to take min and max calculated from the user input
        pickerInput('criteriaX', 'X', choices = criteria, selected = criteria, options=pickerOptions(liveSearch=T)),                 # X coordinates of the plot
        pickerInput('criteriaY', 'Y', choices = criteria, selected = criteria, options=pickerOptions(liveSearch=T))                  # y coordinates of the plot
    ),
    mainPanel(
        fluidRow(
            column(width = 4,
                plotOutput("plot1", height = 700, width = 1000) # the main t-sne plot
            )
        ),
        fluidRow(
            actionButton("SavePlot","Save PNG"),    # the save plot in PNG button
            actionButton("saveSVG","Save SVG")    # the save plot in PNG button
        )
    )
)


########## rendering ###########
server <- function(input, output) {

    #saved_values=reactiveValues()


    # dataset <- reactive({
    #     samples
    # })

    sampleComb  <- reactive({
        if (input$o == "-"){
            norm(samples[,input$z]) - norm(samples[,input$w])
        } else if (input$o == "+") {
            norm(samples[,input$z]) + norm(samples[,input$w])
        } else  {
            norm(samples[,input$z]) * norm(samples[,input$w])
        }

     })

    # slider for color scale
    output$scaleBar <- renderUI({
    lowBond <- min(sampleComb())  # low position of slider
    maxBond <- max(sampleComb())
    sliderInput("colorSelect",
                "Color scale on data score range:",
                min = min(sampleComb()) -  0.3 * max(sampleComb()), #
                max = 1.3 * max(sampleComb()), # 
                value = c(lowBond, maxBond), # 
                )
    })

    output$densityBar <- renderUI({
    min_expression=min(sampleComb())
    med_expression=((input$colorSelect[2]-min_expression)/1.5)+min_expression
    sliderInput("densitySelect",
                "Density scale on data score range:",
                min = min(sampleComb()) -  0.3 * max(sampleComb()), #
                max = 1.3 * max(sampleComb()), # 
                value = med_expression, # 
                )
    })

    
    # density plot
    output$plot2 <- renderPlot({
        plot(density(sampleComb()),  col = "blue", main="", xlab="signature score", ylab="frequency")
    })
    # main plot
    output$plot1 <- renderPlot({
        jet.colors <- colorRampPalette(c("#00007F", "blue", "#007FFF", "cyan", "#7FFF7F", "yellow", "#FF7F00", "red", "#7F0000"))

        max_expression=max(sampleComb())
        min_expression=min(sampleComb())
        med_expression=((input$colorSelect[2]-min_expression)/1.5)+min_expression
        max_tsneX=arrondi(max(samples[,input$criteriaX]))+1
        min_tsneX=arrondi(min(samples[,input$criteriaX]))-1
        max_tsneY=arrondi(max(samples[,input$criteriaY]))+1
        min_tsneY=arrondi(min(samples[,input$criteriaY]))-1
        ggplot(samples , aes_string(x = input$criteriaX, y = input$criteriaY)) +
        ggtitle(paste0(input$z, " vs ", input$w))+
        stat_density_2d(data=samples[sampleComb()>input$densitySelect,], aes(fill=..level..,alpha=..level..),geom='polygon',colour='darkgrey',n=50) +
        geom_jitter(aes(x = samples[,input$criteriaX], y = samples[,input$criteriaY], color=sampleComb()), size=1, alpha=0.7) +
        labs(color = "color scale") + # color legend
        scale_colour_gradientn(colours = jet.colors(10), limits=c(input$colorSelect[1], input$colorSelect[2])) +
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
        ylab(input$criteriaY)
  })


# save plot button
  observeEvent(input$SavePlot,{
    # progress bar
    withProgress(message = 'Saving plot', {
        inputPlot <-     jet.colors <- colorRampPalette(c("#00007F", "blue", "#007FFF", "cyan", "#7FFF7F", "yellow", "#FF7F00", "red", "#7F0000"))

        med_expression=median(sampleComb())
        max_expression=max(sampleComb())
        min_expression=min(sampleComb())
        med_expression=((input$colorSelect[2]-min_expression)/1.5)+min_expression
        max_tsneX=arrondi(max(samples[,input$criteriaX]))+1
        min_tsneX=arrondi(min(samples[,input$criteriaX]))-1
        max_tsneY=arrondi(max(samples[,input$criteriaY]))+1
        min_tsneY=arrondi(min(samples[,input$criteriaY]))-1
        p<-ggplot(samples , aes_string(x = input$criteriaX, y = input$criteriaY)) +
        stat_density_2d(data=samples[sampleComb()>input$densitySelect,], aes(fill=..level..,alpha=..level..),geom='polygon',colour='darkgrey',n=50) +
        geom_jitter(aes(x = samples[,input$criteriaX], y = samples[,input$criteriaY], color=sampleComb()), size=1, alpha=0.7) +
        labs(color = "color scale") + # color legend
        scale_colour_gradientn(colours = jet.colors(10), limits=c(input$colorSelect[1], input$colorSelect[2])) +
        xlab(criteriaX) +
        theme_bw() +
        theme(panel.grid.major = element_blank(), panel.grid.minor = element_blank()) +  # remove grid
        scale_alpha_continuous(range=c(0.1,0.5))+
        scale_x_continuous(limits = c(min_tsneX, max_tsneX))+
        scale_y_continuous(limits = c(min_tsneY, max_tsneY))+
        ggtitle(paste0(input$z, " vs ", input$w))+
        theme(
        plot.title = element_text(size = 14),
        axis.text.x = element_text(size = 14, color='black'),
        axis.text.y = element_text(size = 14, color='black'),
        strip.text.x = element_text(size = 14, color='black',angle=90),
        strip.text.y = element_text(size = 14, color='black',angle=90),
        axis.title.x = element_text(size = 14, color='black'),
        axis.title.y = element_text(size = 14,angle=90))+
        ylab(criteriaY)#+
        #labs(title = input$z)
        file<-paste0("plots/", input$z, "_", input$w, ".png")
        ggsave(file,p, width = 243, height = 150, units = "mm", dpi = 200)
        print("image saved in plots directory !")
    })

  })

    observeEvent(input$saveSVG,{
        # progress bar
    withProgress(message = 'Saving plot', {
        inputPlot <-     jet.colors <- colorRampPalette(c("#00007F", "blue", "#007FFF", "cyan", "#7FFF7F", "yellow", "#FF7F00", "red", "#7F0000"))

        med_expression=median(sampleComb())
        max_expression=max(sampleComb())
        min_expression=min(sampleComb())
        med_expression=((input$colorSelect[2]-min_expression)/1.5)+min_expression
        max_tsneX=arrondi(max(samples[,input$criteriaX]))+1
        min_tsneX=arrondi(min(samples[,input$criteriaX]))-1
        max_tsneY=arrondi(max(samples[,input$criteriaY]))+1
        min_tsneY=arrondi(min(samples[,input$criteriaY]))-1
        p<-ggplot(samples , aes_string(x = input$criteriaX, y = input$criteriaY)) +
        stat_density_2d(data=samples[sampleComb()>input$densitySelect,], aes(fill=..level..,alpha=..level..),geom='polygon',colour='darkgrey',n=50) +
        geom_jitter(aes(x = samples[,input$criteriaX], y = samples[,input$criteriaY], color=sampleComb()), size=1, alpha=0.7) +
        labs(color = "color scale") + # color legend
        scale_colour_gradientn(colours = jet.colors(10), limits=c(input$colorSelect[1], input$colorSelect[2])) +
        xlab(criteriaX) +
        theme_bw() +
        theme(panel.grid.major = element_blank(), panel.grid.minor = element_blank()) +  # remove grid
        scale_alpha_continuous(range=c(0.1,0.5))+
        scale_x_continuous(limits = c(min_tsneX, max_tsneX))+
        scale_y_continuous(limits = c(min_tsneY, max_tsneY))+
        ggtitle(paste0(input$z, " vs ", input$w))+
        theme(
        plot.title = element_text(size = 14),
        axis.text.x = element_text(size = 14, color='black'),
        axis.text.y = element_text(size = 14, color='black'),
        strip.text.x = element_text(size = 14, color='black',angle=90),
        strip.text.y = element_text(size = 14, color='black',angle=90),
        axis.title.x = element_text(size = 14, color='black'),
        axis.title.y = element_text(size = 14,angle=90))+
        ylab(criteriaY)#+
        #labs(title = input$z)
        file<-paste0("plots/", input$z, "_", input$w, ".svg")
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


    save_plot_file <- function(inputPlot){
        file_name <- paste("plots/" , now(), ".svg")
        ggsave(file_name, plot = inputPlot, device = "svg")
    }

}

shinyApp(ui, server)
