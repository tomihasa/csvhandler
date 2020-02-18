package main



import (
    "encoding/csv"
    "fmt"
    "path/filepath"
    "os"
    "time"
    "log"
    "strconv"
    "sort"
)





// lon = x (-73...)   ||   -74... -73... -72
//
// lat = y (40...)    ||   41
//                    ||   40
//                    ||   39
//
// in the CSV file in order
//     lat, lon (y, x)
//
// in the HTML file in order   (HTML file is not part of this project)
//     lon, lat (x, y)
type Station struct {

    stationIdInt  int
    stationName   string
    lon           string
    lat           string
}



// struct for Bike
// using a struct in case of future expansion of code
type Bike struct {

    bikeIdInt  int
}


// time, frequency and average
// used for saving trip times in arrays
// ------------------------------------
//   format of arrays using TimeFrequencyAverage
//       cell 24 = - sum of all trip times (not from cells 00–23)
//                 - frequency of all trip times (not from cells 00–23)
//                 - average of all trip times (not from cells 00–23)
//   cells 00–23 = - average of trip times
//                 - frequency of trip times
//                 - sum of trip times
type TimeFrequencyAverage struct {

    // used when the array has been sorted by time (instead of hour)
    hour       int      // starting hour (0–23) or for sum of them all (24)

    time       int64    // sum of trip time values (seconds)

    frequency  int64    // sum of frequencies of trip times

    average    float64  // average from frequency and total time
}



type Connection struct {

    name       string

    frequency  int
}



const (

    INF01 = "**************************************************"
    INF02 = "* CsvHandler v1.00 (February 18th, 2020)         *"
    INF03 = "*                                                *"
    INF04 = "* Copyright (C) 2020 Tomi Häsä                   *"
    INF05 = "*                                                *"
    INF06 = "* https://github.com/tomihasa/CsvHandler         *"
    INF07 = "*                                                *"
    INF08 = "* To be used with Citi Bike NYC CSV files        *"
    INF09 = "*                                                *"
    INF10 = "* CSV files are available at:                    *"
    INF11 = "* - https://www.citibikenyc.com/system-data      *"
    INF12 = "* - https://s3.amazonaws.com/tripdata/index.html *"
    INF13 = "**************************************************"

    NUMBEROFSTATIONSPERLINE    = 10
    NUMBEROFBIKESPERLINE       = 10
    NUMBEROFCONNECTIONSPERLINE =  4

    GENDERUNKNOWN  = 0
    GENDERMALE     = 1
    GENDERFEMALE   = 2

    DELIMITER = "-"
)



var (

    currentPath       string  =  ""

    stationIdString   string  =  ""
    stationIdInt      int     =  0
    stationName       string  =  ""
    lon               string  =  ""
    lat               string  =  ""

    bikeIdString      string  =  ""
    bikeIdInt         int     =  0

    currentString     string  =  ""

    currentInt        int     =  0

    currentItemNo     int     =  0

    currentStation    = Station{ 0, "", "", "" }

    currentBike       = Bike{ 0 }

    currentHour64     int64 = 0

    currentTime       = TimeFrequencyAverage{ 0, 0, 0, 0}

    durationInt       int     =  0

    duration64        int64   =  0

    temp64            int64   =  0

    startInt          int     =  0

    genderString      string  =  ""

    genderInt         int     =  0

    err               error

    no00              = []byte( "0" )
    no01              = []byte( "1" )
    no02              = []byte( "2" )
    no03              = []byte( "3" )
    no04              = []byte( "4" )
    no05              = []byte( "5" )
    no06              = []byte( "6" )
    no07              = []byte( "7" )
    no08              = []byte( "8" )
    no09              = []byte( "9" )
)



var inf       [13]string

var headings  [15]string

var stationMap    = make( map[ int ]Station )

var stationArray  []Station

var bikeMap       = make( map[ int ]Bike )

var bikeArray     []Bike

var stationMenStartMap      = make( map[ int ]TimeFrequencyAverage )

var stationMenStartArray    []TimeFrequencyAverage

var stationWomenStartMap    = make( map[ int ]TimeFrequencyAverage )

var stationWomenStartArray  []TimeFrequencyAverage

var stationMenEndMap      = make( map[ int ]TimeFrequencyAverage )

var stationMenEndArray    []TimeFrequencyAverage

var stationWomenEndMap    = make( map[ int ]TimeFrequencyAverage )

var stationWomenEndArray  []TimeFrequencyAverage

var connectionMenMap      = make( map[ string ]Connection )

var connectionMenArray    []Connection

var connectionWomenMap    = make( map[ string ]Connection )

var connectionWomenArray  []Connection

var t00  =  TimeFrequencyAverage{  0, 0, 0, 0 }

var t01  =  TimeFrequencyAverage{  1, 0, 0, 0 }

var t02  =  TimeFrequencyAverage{  2, 0, 0, 0 }

var t03  =  TimeFrequencyAverage{  3, 0, 0, 0 }

var t04  =  TimeFrequencyAverage{  4, 0, 0, 0 }

var t05  =  TimeFrequencyAverage{  5, 0, 0, 0 }

var t06  =  TimeFrequencyAverage{  6, 0, 0, 0 }

var t07  =  TimeFrequencyAverage{  7, 0, 0, 0 }

var t08  =  TimeFrequencyAverage{  8, 0, 0, 0 }

var t09  =  TimeFrequencyAverage{  9, 0, 0, 0 }

var t10  =  TimeFrequencyAverage{ 10, 0, 0, 0 }

var t11  =  TimeFrequencyAverage{ 11, 0, 0, 0 }

var t12  =  TimeFrequencyAverage{ 12, 0, 0, 0 }

var t13  =  TimeFrequencyAverage{ 13, 0, 0, 0 }

var t14  =  TimeFrequencyAverage{ 14, 0, 0, 0 }

var t15  =  TimeFrequencyAverage{ 15, 0, 0, 0 }

var t16  =  TimeFrequencyAverage{ 16, 0, 0, 0 }

var t17  =  TimeFrequencyAverage{ 17, 0, 0, 0 }

var t18  =  TimeFrequencyAverage{ 18, 0, 0, 0 }

var t19  =  TimeFrequencyAverage{ 19, 0, 0, 0 }

var t20  =  TimeFrequencyAverage{ 20, 0, 0, 0 }

var t21  =  TimeFrequencyAverage{ 21, 0, 0, 0 }

var t22  =  TimeFrequencyAverage{ 22, 0, 0, 0 }

var t23  =  TimeFrequencyAverage{ 23, 0, 0, 0 }

var t24  =  TimeFrequencyAverage{ 24, 0, 0, 0 }



func init() {

    inf = [13]string { INF01, INF02, INF03, INF04, INF05, INF06, INF07,
          INF08, INF09, INF10, INF11, INF12, INF13 }

    headings = [15]string { "tripduration", "starttime", "stoptime",
               "start station id", "start station name", "start station latitude",
               "start station longitude", "end station id", "end station name",
               "end station latitude", "end station longitude", "bikeid",
               "usertype", "birth year", "gender" }

    stationMenStartMap = map[int]TimeFrequencyAverage {
     0: t00,  1: t01,  2: t02,  3: t03,  4: t04,  5: t05,  6: t06,  7: t07,  8: t08,
     9: t09, 10: t10, 11: t11, 12: t12, 13: t13, 14: t14, 15: t15, 16: t16, 17: t17,
    18: t18, 19: t19, 20: t20, 21: t21, 22: t22, 23: t23, 24: t24, }

    stationWomenStartMap = map[int]TimeFrequencyAverage {
     0: t00,  1: t01,  2: t02,  3: t03,  4: t04,  5: t05,  6: t06,  7: t07,  8: t08,
     9: t09, 10: t10, 11: t11, 12: t12, 13: t13, 14: t14, 15: t15, 16: t16, 17: t17,
    18: t18, 19: t19, 20: t20, 21: t21, 22: t22, 23: t23, 24: t24, }

    stationMenEndMap = map[int]TimeFrequencyAverage {
     0: t00,  1: t01,  2: t02,  3: t03,  4: t04,  5: t05,  6: t06,  7: t07,  8: t08,
     9: t09, 10: t10, 11: t11, 12: t12, 13: t13, 14: t14, 15: t15, 16: t16, 17: t17,
    18: t18, 19: t19, 20: t20, 21: t21, 22: t22, 23: t23, 24: t24, }

    stationWomenEndMap = map[int]TimeFrequencyAverage {
     0: t00,  1: t01,  2: t02,  3: t03,  4: t04,  5: t05,  6: t06,  7: t07,  8: t08,
     9: t09, 10: t10, 11: t11, 12: t12, 13: t13, 14: t14, 15: t15, 16: t16, 17: t17,
    18: t18, 19: t19, 20: t20, 21: t21, 22: t22, 23: t23, 24: t24, }
}



func readFile( filename string ) [][]string {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    fmt.Println( "Current directory: " + startpath )
    fmt.Println()

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println( "Attempting to open file: " + filename )

    os.Chdir( currentPath )

    // avaa .csv-tiedosto

    file, err := os.Open( currentPath )

    if err != nil {

        log.Fatalf( "Opening file failed: %s", err )
    }

    defer file.Close()

    // file size

    fileInfo, err := file.Stat()

    if err != nil {

        log.Fatalf( "File size could not be requested from system: %s", err )
    }

    fmt.Println()
    fmt.Printf( "File size: %d", fileInfo.Size() )
    fmt.Println( " bytes" )
    fmt.Println()

    fmt.Println( "Processing data. Might take a few seconds..." )
    fmt.Println()

    // lue .csv-tiedosto

    csvReader := csv.NewReader( file )

    csvData, err := csvReader.ReadAll()

    if err != nil {

        log.Fatalf( "Reading file failed: %s", err )
    }

    // read all rows except the heading row
    csvDataWithoutHeadings := csvData[1:]

    fmt.Printf( "Number of rows in file: %d", len( csvData ) )
    fmt.Println()

    return csvDataWithoutHeadings
}



func writeStationCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range stationArray {

        stationId    := strconv.Itoa( row.stationIdInt )

        stationName  := "\"" + row.stationName + "\""

        lon          := row.lon

        lat          := row.lat

        array        := createArrayForCsv4( stationId, stationName, lon, lat )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func writeBikeCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range bikeArray {

        bikeId  := strconv.Itoa( row.bikeIdInt )

        array   := createArrayForCsv1( bikeId )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func writeStartStationMenCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range stationMenStartArray {

        hour       := strconv.Itoa( row.hour )

        time       := strconv.FormatInt( row.time, 10 )

        frequency  := strconv.FormatInt( row.frequency, 10 )

        average    := strconv.FormatFloat( row.average, 'f', 2, 64 )

        array      := createArrayForCsv4( hour, time, frequency, average )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}


func writeStartStationWomenCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range stationWomenStartArray {

        hour       := strconv.Itoa( row.hour )

        time       := strconv.FormatInt( row.time, 10 )

        frequency  := strconv.FormatInt( row.frequency, 10 )

        average    := strconv.FormatFloat( row.average, 'f', 2, 64 )

        array      := createArrayForCsv4( hour, time, frequency, average )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func writeEndStationMenCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range stationMenEndArray {

        hour       := strconv.Itoa( row.hour )

        time       := strconv.FormatInt( row.time, 10 )

        frequency  := strconv.FormatInt( row.frequency, 10 )

        average    := strconv.FormatFloat( row.average, 'f', 2, 64 )

        array      := createArrayForCsv4( hour, time, frequency, average )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func writeEndStationWomenCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range stationWomenEndArray {

        hour       := strconv.Itoa( row.hour )

        time       := strconv.FormatInt( row.time, 10 )

        frequency  := strconv.FormatInt( row.frequency, 10 )

        average    := strconv.FormatFloat( row.average, 'f', 2, 64 )

        array      := createArrayForCsv4( hour, time, frequency, average )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func writeConnectionMenCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range connectionMenArray {

        name       := "\"" + row.name + "\""

        frequency  := strconv.Itoa( row.frequency )

        array      := createArrayForCsv2( name, frequency )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func writeConnectionWomenCsv( filename string ) {

    // path

    startpath, errd := filepath.Abs( filepath.Dir( os.Args[0] ) )

    if errd != nil {

        log.Fatal( errd )
    }

    currentPath = filepath.Join( startpath, filename )

    fmt.Println()
    fmt.Println( "Writing file: "+ filename )

    // create file

    csvFile, err := os.Create( currentPath )

    if err != nil {

        log.Fatalf( "Creating file failed: %s", err )
    }

    // write to file

    csvWriter := csv.NewWriter( csvFile )
 
    for _, row := range connectionWomenArray {

        name       := "\"" + row.name + "\""

        frequency  := strconv.Itoa( row.frequency )

        array      := createArrayForCsv2( name, frequency )

        _ = csvWriter.Write( array )
    }

    csvWriter.Flush()

    csvFile.Close()
}



func createArrayForCsv1( x string ) []string {

    var array []string

    array = append( array, x )

    return  array
}



func createArrayForCsv2( x1 string, x2 string ) []string {

    var array []string

    array = append( array, x1, x2 )

    return  array
}



func createArrayForCsv4( x1 string, x2 string, x3 string, x4 string ) []string {

    var array []string

    array = append( array, x1, x2, x3, x4 )

    return  array
}



// brute force
//   strconv.ParseInt and strconv.Atoi do not seem to work with Citi Bike CSV
//   files when converting from string to int, so lets use brute force instead
func stringToIntBruteForce( characters string ) int {

    var number = 0

    if ( characters[0] == no00[0] && characters[1] == no00[0] ) {

         number = 0
    }

    if ( characters[0] == no00[0] && characters[1] == no01[0] ) {

         number = 1
    }

    if ( characters[0] == no00[0] && characters[1] == no02[0] ) {

        number = 2
    }

    if ( characters[0] == no00[0] && characters[1] == no03[0] ) {

        number = 3
    }

    if ( characters[0] == no00[0] && characters[1] == no04[0] ) {

        number = 4
    }

    if ( characters[0] == no00[0] && characters[1] == no05[0] ) {

        number = 5
    }

    if ( characters[0] == no00[0] && characters[1] == no06[0] ) {

        number = 6
    }

    if ( characters[0] == no00[0] && characters[1] == no07[0] ) {

        number = 7
    }

    if ( characters[0] == no00[0] && characters[1] == no08[0] ) {

        number = 8
    }

    if ( characters[0] == no00[0] && characters[1] == no09[0] ) {

        number = 9
    }

    if ( characters[0] == no01[0] && characters[1] == no00[0] ) {

        number = 10
    }

    if ( characters[0] == no01[0] && characters[1] == no01[0] ) {

        number = 11
    }

    if ( characters[0] == no01[0] && characters[1] == no02[0] ) {

        number = 12
    }

    if ( characters[0] == no01[0] && characters[1] == no03[0] ) {

        number = 13
    }

    if ( characters[0] == no01[0] && characters[1] == no04[0] ) {

        number = 14
    }

    if ( characters[0] == no01[0] && characters[1] == no05[0] ) {

        number = 15
    }

    if ( characters[0] == no01[0] && characters[1] == no06[0] ) {

        number = 16
    }

    if ( characters[0] == no01[0] && characters[1] == no07[0] ) {

        number = 17
    }

    if ( characters[0] == no01[0] && characters[1] == no08[0] ) {

        number = 18
    }

    if ( characters[0] == no01[0] && characters[1] == no09[0] ) {

        number = 19
    }

    if ( characters[0] == no02[0] && characters[1] == no00[0] ) {

        number = 20
    }

    if ( characters[0] == no02[0] && characters[1] == no01[0] ) {

        number = 21
    }

    if ( characters[0] == no02[0] && characters[1] == no02[0] ) {

        number = 22
    }

    if ( characters[0] == no02[0] && characters[1] == no03[0] ) {

        number = 23
    }

        return number
}



// --------- START STATIONS: AVERAGE ---------

// start stations: average times for men

func countAverageStartMen() {

    for i := 0; i < 25; i++ {

        timeMen := stationMenStartArray[i].time

        frequencyMen := stationMenStartArray[i].frequency

        if ( frequencyMen != 0 ) {

            stationMenStartArray[i].average = float64( timeMen ) / float64( frequencyMen )
        }
    }
}


// start stations: average times for women

func countAverageStartWomen() {

    for j := 0; j < 25; j++ {

        timeWomen := stationWomenStartArray[j].time

        frequencyWomen := stationWomenStartArray[j].frequency

        if ( frequencyWomen != 0 ) {

            stationWomenStartArray[j].average = float64( timeWomen ) / float64( frequencyWomen )
        }
    }
}



// --------- END STATIONS: AVERAGE ---------

// end stations: average times for men

func countAverageEndMen() {

    for i := 0; i < 25; i++ {

        timeMen := stationMenEndArray[i].time

        frequencyMen := stationMenEndArray[i].frequency

        if ( frequencyMen != 0 ) {

            stationMenEndArray[i].average = float64( timeMen ) / float64( frequencyMen )
        }
    }
}


// end stations: average times for women

func countAverageEndWomen() {

    for j := 0; j < 25; j++ {

        timeWomen := stationWomenEndArray[j].time

        frequencyWomen := stationWomenEndArray[j].frequency

        if ( frequencyWomen != 0 ) {

            stationWomenEndArray[j].average = float64( timeWomen ) / float64( frequencyWomen )
        }
    }
}






func main() {

    // --------- PROGRAM INFO ---------

    timeWhenProgramStarted := time.Now()

    for _, programinfoline := range inf {

        fmt.Println( programinfoline )
    }
    fmt.Println()



    // --------- READ CSV FILE ---------

    tripdata := readFile( "citibike.csv" )
    fmt.Println()


    // --------- STATIONS ---------

    // start stations to map
    for _, row := range tripdata {

        stationIdString  = row[3]
        stationName      = row[4]
        lat              = row[5]
        lon              = row[6]
        stationIdInt, _  = strconv.Atoi( stationIdString )

        currentStation = Station{ stationIdInt, stationName, lon, lat }

        stationMap[ stationIdInt ] = currentStation
    }

    // end stations to map
    for _, row := range tripdata {

        stationIdString  = row[7]
        stationName      = row[8]
        lat              = row[9]
        lon              = row[10]
        stationIdInt, _  = strconv.Atoi( stationIdString )

        currentStation = Station{ stationIdInt, stationName, lon, lat }

        stationMap[ stationIdInt ] = currentStation
    }


    // save data from map to array for sorting
    for currentString, _ := range stationMap {

        currentStation = stationMap[ currentString ]

        stationArray = append( stationArray, currentStation )
    }


    // sort by station id
    sort.Slice( stationArray, func(i, j int) bool { return stationArray[i].stationIdInt < stationArray[j].stationIdInt } )


    fmt.Printf( "%d stations: ", len( stationMap ) )
    fmt.Println()

    for _, line := range stationArray {

        fmt.Println( line )
    }






    // --------- BIKES ---------

    // bikes to map
    for _, row := range tripdata {

        bikeIdString = row[11]
        bikeIdInt, _ = strconv.Atoi( bikeIdString )

        currentBike = Bike{ bikeIdInt }

        bikeMap[ bikeIdInt ] = currentBike
    }


    // save data from map to array for sorting
    for currentInt, _ := range bikeMap {

        currentBike = bikeMap[ currentInt ]

        bikeArray = append( bikeArray, currentBike )
    }


    // sort by bike id
    sort.Slice( bikeArray, func(i, j int) bool { return bikeArray[i].bikeIdInt < bikeArray[j].bikeIdInt } )


    fmt.Println()
    fmt.Printf( "%d bikes: ", len( bikeMap ) )
    fmt.Println()

    currentItemNo = 0

    for _, line := range bikeArray {

        currentItemNo++

        bikeInArray := line.bikeIdInt

        fmt.Printf( " %7d ", bikeInArray )

        if ( currentItemNo >= NUMBEROFBIKESPERLINE ) {

            println()
            currentItemNo = 0
        }
    }
    fmt.Println()






    // --------- START STATIONS ---------

    for _, row := range tripdata {

        durationString   := row[0]
        duration64, _    = strconv.ParseInt( durationString, 10, 64 )

        startTimeString      := row[1]
        startTimeSlice       := startTimeString[11:13]
        startTimeCharacters  := string( startTimeSlice )

        hour := stringToIntBruteForce( startTimeCharacters )

        genderString    = row[14]
        genderInt, _    = strconv.Atoi( genderString )

        if ( genderInt == GENDERMALE ) {

            currentMen := stationMenStartMap[ hour ]

            currentTimesMen := currentMen.time + duration64

            currentMen.time = currentTimesMen

            currentMen.frequency = currentMen.frequency + 1

            stationMenStartMap[ hour ] = currentMen

            // sum of all times and frequencies
            sumOfMen := stationMenStartMap[ 24 ]

            currentOfAllTripTimesMen := sumOfMen.time + duration64

            sumOfMen.time = currentOfAllTripTimesMen

            sumOfMen.frequency = sumOfMen.frequency + 1

            stationMenStartMap[ 24 ] = sumOfMen
        }

        if ( genderInt == GENDERFEMALE ) {

            currentWomen := stationWomenStartMap[ hour ]

            currentTimesWomen := currentWomen.time + duration64

            currentWomen.time = currentTimesWomen

            currentWomen.frequency = currentWomen.frequency + 1

            stationWomenStartMap[ hour ] = currentWomen

            // sum of all times and frequencies
            sumOfWomen := stationWomenStartMap[ 24 ]

            currentOfAllTripTimesWomen := sumOfWomen.time + duration64

            sumOfWomen.time = currentOfAllTripTimesWomen

            sumOfWomen.frequency = sumOfWomen.frequency + 1

            stationWomenStartMap[ 24 ] = sumOfWomen
        }
    }



    // --------- START STATIONS: MEN ---------

    // save data from map to array for sorting
    for i := 0; i < 25; i++ {

        currentTime = stationMenStartMap[ i ]

        stationMenStartArray = append( stationMenStartArray, currentTime )
    }


    // sort by time
    sort.Slice( stationMenStartArray, func(i, j int) bool { return stationMenStartArray[i].time > stationMenStartArray[j].time } )

    countAverageStartMen()

    fmt.Println()
    fmt.Println( "Start stations, men (starting hour, sum of trip durations, number of trips, average duration):" )
    fmt.Println( "hour          sum  frequency     average" )

    for i := 0; i < 25; i++ {

        currentTime = stationMenStartArray[ i ]

        fmt.Printf( "%2d %14d %10d  %10.2f", currentTime.hour, currentTime.time, currentTime.frequency, currentTime.average )
        fmt.Println()
    }



    // --------- START STATIONS: WOMEN ---------

    // save data from map to array for sorting
    for i := 0; i < 25; i++ {

        currentTime = stationWomenStartMap[ i ]

        stationWomenStartArray = append( stationWomenStartArray, currentTime )
    }


    // sort by time
    sort.Slice( stationWomenStartArray, func(i, j int) bool { return stationWomenStartArray[i].time > stationWomenStartArray[j].time } )

    countAverageStartWomen()

    fmt.Println()
    fmt.Println( "Start stations, women (starting hour, sum of trip durations, number of trips, average duration):" )
    fmt.Println( "hour          sum  frequency     average" )

    for i := 0; i < 25; i++ {

        currentTime = stationWomenStartArray[ i ]

        fmt.Printf( "%2d %14d %10d  %10.2f", currentTime.hour, currentTime.time, currentTime.frequency, currentTime.average )
        fmt.Println()
    }






    // --------- END STATIONS ---------

    for _, row := range tripdata {

        durationString   := row[0]
        duration64, _    = strconv.ParseInt( durationString, 10, 64 )

        endTimeString      := row[2]
        endTimeSlice       := endTimeString[11:13]
        endTimeCharacters  := string( endTimeSlice )

        hour := stringToIntBruteForce( endTimeCharacters )

        genderString    = row[14]
        genderInt, _    = strconv.Atoi( genderString )

        if ( genderInt == GENDERMALE ) {

            currentMen := stationMenEndMap[ hour ]

            currentTimesMen := currentMen.time + duration64

            currentMen.time = currentTimesMen

            currentMen.frequency = currentMen.frequency + 1

            stationMenEndMap[ hour ] = currentMen

            // sum of all times and frequencies
            sumOfMen := stationMenEndMap[ 24 ]

            currentOfAllTripTimesMen := sumOfMen.time + duration64

            sumOfMen.time = currentOfAllTripTimesMen

            sumOfMen.frequency = sumOfMen.frequency + 1

            stationMenEndMap[ 24 ] = sumOfMen
        }

        if ( genderInt == GENDERFEMALE ) {

            currentWomen := stationWomenEndMap[ hour ]

            currentTimesWomen := currentWomen.time + duration64

            currentWomen.time = currentTimesWomen

            currentWomen.frequency = currentWomen.frequency + 1

            stationWomenEndMap[ hour ] = currentWomen

            // sum of all times and frequencies
            sumOfWomen := stationWomenEndMap[ 24 ]

            currentOfAllTripTimesWomen := sumOfWomen.time + duration64

            sumOfWomen.time = currentOfAllTripTimesWomen

            sumOfWomen.frequency = sumOfWomen.frequency + 1

            stationWomenEndMap[ 24 ] = sumOfWomen
        }
    }


    // --------- END STATIONS: MEN ---------

    // save data from map to array for sorting
    for i := 0; i < 25; i++ {

        currentTime = stationMenEndMap[ i ]

        stationMenEndArray = append( stationMenEndArray, currentTime )
    }


    // sort by time
    sort.Slice( stationMenEndArray, func(i, j int) bool { return stationMenEndArray[i].time > stationMenEndArray[j].time } )

    countAverageEndMen()

    fmt.Println()
    fmt.Println( "End stations, men (starting hour, sum of trip durations, number of trips, average duration):" )
    fmt.Println( "hour          sum  frequency     average" )

    for i := 0; i < 25; i++ {

        currentTime = stationMenEndArray[ i ]

        fmt.Printf( "%2d %14d %10d  %10.2f", currentTime.hour, currentTime.time, currentTime.frequency, currentTime.average )
        fmt.Println()
    }



    // --------- END STATIONS: WOMEN ---------

    // save data from map to array for sorting
    for i := 0; i < 25; i++ {

        currentTime = stationWomenEndMap[ i ]

        stationWomenEndArray = append( stationWomenEndArray, currentTime )
    }


    // sort by time
    sort.Slice( stationWomenEndArray, func(i, j int) bool { return stationWomenEndArray[i].time > stationWomenEndArray[j].time } )

    countAverageEndWomen()

    fmt.Println()
    fmt.Println( "End stations, women (starting hour, sum of trip durations, number of trips, average duration):" )
    fmt.Println( "hour          sum  frequency     average" )

    for i := 0; i < 25; i++ {

        currentTime = stationWomenEndArray[ i ]

        fmt.Printf( "%2d %14d %10d  %10.2f", currentTime.hour, currentTime.time, currentTime.frequency, currentTime.average )
        fmt.Println()
    }






    // --------- CONNECTIONS: MEN ---------

    // fill map with unique connection ids before counting frequencies
    for _, row := range tripdata {

        connectionStartString := row[3]
        connectionEndString   := row[7]

        genderString  = row[14]
        genderInt, _  = strconv.Atoi( genderString )

        connectionString := connectionStartString + DELIMITER + connectionEndString

        if ( genderInt == GENDERMALE ) {

            connectionMenMap[ connectionString ] = Connection{ connectionString, 0 }
        }
    }


    // count frequencies
    for _, row := range tripdata {

        connectionStartString := row[3]
        connectionEndString   := row[7]

        genderString  = row[14]
        genderInt, _  = strconv.Atoi( genderString )

        connectionString := connectionStartString + DELIMITER + connectionEndString

        if ( genderInt == GENDERMALE ) {

            currentConnection := connectionMenMap[ connectionString ]

            currentConnection.frequency = currentConnection.frequency + 1

            connectionMenMap[ connectionString ] = currentConnection
        }
    }


    // save data from map to array for sorting
    for currentInt, _ := range connectionMenMap {

        currentConnection := connectionMenMap[ currentInt ]

        connectionMenArray = append( connectionMenArray, currentConnection )
    }


    // sort by frequency
    sort.Slice( connectionMenArray, func(i, j int) bool { return connectionMenArray[i].frequency > connectionMenArray[j]. frequency } )

    fmt.Println()
    fmt.Println( "Connections, men:" )

    currentItemNo = 0

    for _, currentConnection := range connectionMenArray {

        currentItemNo++

        fmt.Printf( "%12v %6d   ", currentConnection.name, currentConnection.frequency )

        if ( currentItemNo >= NUMBEROFCONNECTIONSPERLINE ) {

            println()
            currentItemNo = 0
        }

    }
    fmt.Println()



    // --------- CONNECTIONS: WOMEN ---------

    // fill map with unique connection ids before counting frequencies
    for _, row := range tripdata {

        connectionStartString := row[3]
        connectionEndString   := row[7]

        genderString  = row[14]
        genderInt, _  = strconv.Atoi( genderString )

        connectionString := connectionStartString + DELIMITER + connectionEndString

        if ( genderInt == GENDERFEMALE ) {

            connectionWomenMap[ connectionString ] = Connection{ connectionString, 0 }
        }
    }


    // count frequencies
    for _, row := range tripdata {

        connectionStartString := row[3]
        connectionEndString   := row[7]

        genderString  = row[14]
        genderInt, _  = strconv.Atoi( genderString )

        connectionString := connectionStartString + DELIMITER + connectionEndString

        if ( genderInt == GENDERFEMALE ) {

            currentConnection := connectionWomenMap[ connectionString ]

            currentConnection.frequency = currentConnection.frequency + 1

            connectionWomenMap[ connectionString ] = currentConnection
        }
    }


    // save data from map to array for sorting
    for currentInt, _ := range connectionWomenMap {

        currentConnection := connectionWomenMap[ currentInt ]

        connectionWomenArray = append( connectionWomenArray, currentConnection )
    }



    // sort by frequency
    sort.Slice( connectionWomenArray, func(i, j int) bool { return connectionWomenArray[i].frequency > connectionWomenArray[j]. frequency } )

    fmt.Println()
    fmt.Println( "Connections, women:" )

    currentItemNo = 0

    for _, currentConnection := range connectionWomenArray {

        currentItemNo++

        fmt.Printf( "%12v %6d   ", currentConnection.name, currentConnection.frequency )

        if ( currentItemNo >= NUMBEROFCONNECTIONSPERLINE ) {

            println()
            currentItemNo = 0
        }

    }
    fmt.Println()



    // --------- WRITE CSV FILES ---------

    writeStationCsv( "stations.csv" )

    writeBikeCsv( "bikes.csv" )

    writeStartStationMenCsv( "hoursstartstationmen.csv" )

    writeStartStationWomenCsv( "hoursstartstationwomen.csv" )

    writeEndStationMenCsv( "hoursendstationmen.csv" )

    writeEndStationWomenCsv( "hoursendstationwomen.csv" )

    writeConnectionMenCsv( "connectionsmen.csv" )

    writeConnectionWomenCsv( "connectionswomen.csv" )



    // --------- TIME TO EXECUTE PROGRAM ---------

    timeToExecuteProgram := time.Since( timeWhenProgramStarted )

    fmt.Println()
    fmt.Printf( "Time taken to execute program: %s", timeToExecuteProgram )
    fmt.Println()
}

