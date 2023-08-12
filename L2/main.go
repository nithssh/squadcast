package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// file io helper. https://gobyexample.com/reading-files
func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Movie struct {
	id           int64
	title        string
	year         int64
	country      string
	genere       string
	director     string
	mins         int64
	poster       string
	avg_rating   float64
	rating_count int64
}

type Rater struct {
	avg_rating   float64
	rating_count int64
}

func filterMovies(m []Movie, condition func(Movie) bool) []Movie {
	filteredSlice := make([]Movie, 0)

	for _, element := range m {
		if condition(element) {
			filteredSlice = append(filteredSlice, element)
		}
	}

	return filteredSlice
}

func main() {
	movFile, err := os.Open("movies.csv")
	check(err)
	scanner := bufio.NewScanner(movFile)
	scanner.Scan() // discard first line
	movies := make([]Movie, 0)
	for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		fields, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}

		f0, _ := strconv.ParseInt(fields[0], 10, 64)
		f2, _ := strconv.ParseInt(fields[2], 10, 64)
		f6, _ := strconv.ParseInt(fields[6], 10, 64)
		movies = append(
			movies,
			Movie{f0, fields[1], f2, fields[3], fields[4], fields[5], f6, fields[7], 0, 0})
	}

	fmt.Println("\n1.a Five longest movies: ")
	sort.Slice(movies, func(i, j int) bool {
		if movies[i].mins == movies[j].mins {
			return strings.Compare(movies[i].title, movies[j].title) < 0
		}
		return movies[i].mins > movies[j].mins
	})
	for i := 0; i < 5; i++ {
		fmt.Println(movies[i].title+"; Duration In Minutes:", movies[i].mins)
	}

	fmt.Println("\n1.b Five latest movies: ")
	sort.Slice(movies, func(i, j int) bool {
		if movies[i].year == movies[j].year {
			return strings.Compare(movies[i].title, movies[j].title) < 0
		}
		return movies[i].year > movies[j].year
	})
	for i := 0; i < 5; i++ {
		fmt.Println(movies[i].title+"; Year:", movies[i].year)
	}

	// sort by movie id at last
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].id < movies[j].id
	})

	// for _, m := range movies {
	// 	fmt.Println("mid:", m.id)
	// }

	ratFile, err := os.Open("ratings.csv")
	check(err)
	scanner2 := bufio.NewScanner(ratFile)
	scanner2.Scan() // discard first line
	raters := make(map[string]Rater, 1000)

	// process all the rating
	for scanner2.Scan() {
		line := scanner2.Text()
		reader2 := csv.NewReader(strings.NewReader(line))
		fields, err := reader2.Read()
		if err != nil {
			log.Fatal(err)
		}

		mid, _ := strconv.ParseInt(fields[1], 10, 64)
		rating, _ := strconv.ParseFloat(fields[2], 64)

		// update the individual rater
		if rater, ok := raters[fields[0]]; ok {
			rater.avg_rating = ((float64(rater.rating_count) * rater.avg_rating) + rating) / float64(rater.rating_count+1)
			rater.rating_count += 1

			raters[fields[0]] = rater
		} else {
			nrater := Rater{rating, 1}
			raters[fields[0]] = nrater
		}

		// find the movie in slice by id
		// index := sort.Search(len(movies), func(i int) bool {
		// 	return movies[i].id == mid
		// })

		// linear search...
		index := -1
		for i, m := range movies {
			if m.id == mid {
				index = i
			}
		}

		if index >= len(movies) || index == -1 {
			// movie not in DB
			continue
		}

		// update movie's rating record
		movies[index].avg_rating = ((float64(movies[index].rating_count) * movies[index].avg_rating) + rating) / float64(movies[index].rating_count+1)
		movies[index].rating_count += 1
	}

	movies_filt := filterMovies(movies, func(m Movie) bool {
		return m.rating_count >= 5
	})

	//  for _, m := range movies_filt { //
	// 	fmt.Println(m.title, m.rating_count, m.avg_rating)
	// }

	fmt.Println("\n1.c Highest rated movies: ")
	sort.Slice(movies_filt, func(i, j int) bool {
		if movies_filt[i].avg_rating == movies_filt[j].avg_rating {
			return strings.Compare(movies[i].title, movies[j].title) < 0
		}
		return movies_filt[i].avg_rating > movies_filt[j].avg_rating
	})
	for i := 0; i < 5; i++ {
		fmt.Printf("%s; Average Rating: %v\n", movies_filt[i].title, movies_filt[i].avg_rating)
	}

	fmt.Println("\n1.d Most rated movies: ")
	sort.Slice(movies, func(i, j int) bool {
		if movies[i].rating_count == movies[j].rating_count {
			return strings.Compare(movies[i].title, movies[j].title) < 0
		}
		return movies[i].rating_count > movies[j].rating_count
	})
	for i := 0; i < 5; i++ {
		fmt.Printf("%s; Number of Ratings: %v\n", movies[i].title, movies[i].rating_count)
	}

	fmt.Println("\n2. Unique Raters:", len(raters)) 
}
