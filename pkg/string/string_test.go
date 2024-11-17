package string

import (
	"github.com/Confidenceman02/scion-tools/pkg/basics"
	"github.com/Confidenceman02/scion-tools/pkg/char"
	"github.com/Confidenceman02/scion-tools/pkg/list"
	"github.com/Confidenceman02/scion-tools/pkg/maybe"
	"github.com/Confidenceman02/scion-tools/pkg/tuple"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	asserts := assert.New(t)

	t.Run("IsEmpty", func(t *testing.T) {
		SUT1 := IsEmpty("")
		SUT2 := IsEmpty("the world")

		asserts.True(SUT1)
		asserts.False(SUT2)
	})

	t.Run("Length", func(t *testing.T) {
		SUT1 := Length("innumerable")
		SUT2 := Length("")

		asserts.Equal(basics.Int(11), SUT1)
		asserts.Equal(basics.Int(0), SUT2)
	})

	t.Run("Length - Multibyte string", func(t *testing.T) {
		SUT := Length("⌘")

		asserts.Equal(basics.Int(1), SUT)
	})

	t.Run("Reverse", func(t *testing.T) {
		SUT1 := Reverse("stressed")
		SUT2 := Reverse("stressed⌘")

		asserts.Equal(String("desserts"), SUT1)
		asserts.Equal(String("⌘desserts"), SUT2)
	})

	t.Run("Repeat", func(t *testing.T) {
		SUT := Repeat(3, "ha")

		asserts.Equal(String("hahaha"), SUT)
	})

	t.Run("Replace", func(t *testing.T) {
		SUT1 := Replace(".", "-", "Json.Decode.succeed")
		SUT2 := Replace(",", "/", "a,b,c,d")

		asserts.Equal(String("Json-Decode-succeed"), SUT1)
		asserts.Equal(String("a/b/c/d"), SUT2)
	})
}

func TestBuildingAndSPlitting(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Append", func(t *testing.T) {
		SUT := Append("butter", "fly")

		asserts.Equal(String("butterfly"), SUT)
	})

	t.Run("Concat", func(t *testing.T) {
		SUT := Concat(list.FromSlice([]String{"never", "the", "less"}))

		asserts.Equal(String("nevertheless"), SUT)
	})

	t.Run("Split", func(t *testing.T) {
		SUT1 := Split(",", "cat,dog,cow")
		SUT2 := Split("/", "home/evan/Desktop")

		asserts.Equal([]String{"cat", "dog", "cow"}, list.ToSlice(SUT1))
		asserts.Equal([]String{"home", "evan", "Desktop"}, list.ToSlice(SUT2))
	})

	t.Run("Words", func(t *testing.T) {
		SUT := Words("How are \t you? \n Good?")

		asserts.Equal([]String{"How", "are", "you?", "Good?"}, list.ToSlice(SUT))
	})

	t.Run("Lines", func(t *testing.T) {
		SUT := Lines("How are you?\nGood?")

		asserts.Equal([]String{"How are you?", "Good?"}, list.ToSlice(SUT))
	})
}

func TestGetSubscrings(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Slice", func(t *testing.T) {
		SUT1 := Slice(7, 9, "snakes on a plane!")
		SUT2 := Slice(0, 6, "snakes on a plane!")
		SUT3 := Slice(0, -7, "snakes on a plane!")
		SUT4 := Slice(-6, -1, "snakes on a plane!")
		SUT5 := Slice(1, 1, "snakes on a plane!")
		SUT6 := Slice(0, 500, "snakes on a plane!")
		SUT7 := Slice(0, -500, "snakes on a plane!")

		asserts.Equal(String("on"), SUT1)
		asserts.Equal(String("snakes"), SUT2)
		asserts.Equal(String("snakes on a"), SUT3)
		asserts.Equal(String("plane"), SUT4)
		asserts.Equal(String(""), SUT5)
		asserts.Equal(String("snakes on a plane!"), SUT6)
		asserts.Equal(String(""), SUT7)
	})

	t.Run("Left", func(t *testing.T) {
		SUT := Left(2, "Mulder")

		asserts.Equal(String("Mu"), SUT)
	})

	t.Run("Right", func(t *testing.T) {
		SUT := Right(2, "Scully")

		asserts.Equal(String("ly"), SUT)
	})

	t.Run("DropLeft", func(t *testing.T) {
		SUT := DropLeft(2, "The Lone Gunmen")

		asserts.Equal(String("e Lone Gunmen"), SUT)
	})

	t.Run("DropRight", func(t *testing.T) {
		SUT := DropRight(2, "Cigarette Smoking Man")

		asserts.Equal(String("Cigarette Smoking M"), SUT)
	})
}

func TestCheckForSubstrings(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Contains", func(t *testing.T) {
		SUT1 := Contains("the", "theory")
		SUT2 := Contains("hat", "theory")
		SUT3 := Contains("THE", "theory")

		asserts.Equal(true, SUT1)
		asserts.Equal(false, SUT2)
		asserts.Equal(false, SUT3)
	})

	t.Run("StartsWith", func(t *testing.T) {
		SUT1 := StartsWith("the", "theory")
		SUT2 := StartsWith("ory", "theory")

		asserts.Equal(true, SUT1)
		asserts.Equal(false, SUT2)
	})

	t.Run("EndsWith", func(t *testing.T) {
		SUT1 := EndsWith("the", "theory")
		SUT2 := EndsWith("ory", "theory")

		asserts.Equal(false, SUT1)
		asserts.Equal(true, SUT2)
	})

	t.Run("Indexes", func(t *testing.T) {
		SUT1 := Indexes("i", "Mississippi")
		SUT2 := Indexes("ss", "Mississippi")
		SUT3 := Indexes("needle", "haysack")

		asserts.Equal([]basics.Int{1, 4, 7, 10}, list.ToSlice(SUT1))
		asserts.Equal([]basics.Int{2, 5}, list.ToSlice(SUT2))
		asserts.Equal([]basics.Int{}, list.ToSlice(SUT3))
	})
}

func TestIntConversions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ToInt", func(t *testing.T) {
		SUT1 := ToInt("123")
		SUT2 := ToInt("-42")
		SUT3 := ToInt("3.1")
		SUT4 := ToInt("31a")

		asserts.Equal(maybe.Just[basics.Int]{Value: basics.Int(123)}, SUT1)
		asserts.Equal(maybe.Just[basics.Int]{Value: basics.Int(-42)}, SUT2)
		asserts.Equal(maybe.Nothing{}, SUT3)
		asserts.Equal(maybe.Nothing{}, SUT4)
	})

	t.Run("FromInt", func(t *testing.T) {
		SUT1 := FromInt(123)
		SUT2 := FromInt(-42)

		asserts.Equal(String("123"), SUT1)
		asserts.Equal(String("-42"), SUT2)
	})

	t.Run("ToFloat", func(t *testing.T) {
		SUT1 := ToFloat("123")
		SUT2 := ToFloat("-42")
		SUT3 := ToFloat("3.1")
		SUT4 := ToFloat("31a")

		asserts.Equal(maybe.Just[basics.Float]{Value: 123.0}, SUT1)
		asserts.Equal(maybe.Just[basics.Float]{Value: -42.0}, SUT2)
		asserts.Equal(maybe.Just[basics.Float]{Value: 3.1}, SUT3)
		asserts.Equal(maybe.Nothing{}, SUT4)
	})

	t.Run("FromFloat", func(t *testing.T) {
		SUT1 := FromFloat(123)
		SUT2 := FromFloat(-42)
		SUT3 := FromFloat(3.1)

		asserts.Equal(String("123"), SUT1)
		asserts.Equal(String("-42"), SUT2)
		asserts.Equal(String("3.1"), SUT3)
	})
}

func TestCharConversions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("FromChar", func(t *testing.T) {
		SUT := FromChar('a')

		asserts.Equal(String("a"), SUT)
	})

	t.Run("Cons", func(t *testing.T) {
		SUT := Cons('T', "he truth is out there")

		asserts.Equal(String("The truth is out there"), SUT)
	})

	t.Run("Uncons", func(t *testing.T) {
		SUT := Uncons("abc")

		asserts.Equal(maybe.Just[tuple.Tuple2[char.Char, String]]{Value: tuple.Pair(char.Char('a'), String("bc"))}, SUT)
	})
}

func TestListConversions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ToList", func(t *testing.T) {
		SUT1 := ToList("abc")
		SUT2 := ToList("🙈🙉🙊")

		asserts.Equal([]char.Char{'a', 'b', 'c'}, list.ToSlice(SUT1))
		asserts.Equal([]char.Char{'🙈', '🙉', '🙊'}, list.ToSlice(SUT2))
	})

	t.Run("FromList", func(t *testing.T) {
		SUT1 := FromList(list.FromSlice([]char.Char{'a', 'b', 'c'}))
		SUT2 := FromList(list.FromSlice([]char.Char{'🙈', '🙉', '🙊'}))

		asserts.Equal(String("abc"), SUT1)
		asserts.Equal(String("🙈🙉🙊"), SUT2)
	})

}

func TestFormatting(t *testing.T) {
	asserts := assert.New(t)

	t.Run("ToUpper", func(t *testing.T) {
		SUT := ToUpper("skinner")

		asserts.Equal(String("SKINNER"), SUT)
	})
	t.Run("ToLower", func(t *testing.T) {
		SUT := ToLower("X-FILES")

		asserts.Equal(String("x-files"), SUT)
	})

	t.Run("Pad", func(t *testing.T) {
		SUT1 := Pad(5, ' ', "1")
		SUT2 := Pad(5, ' ', "11")
		SUT3 := Pad(5, ' ', "121")

		asserts.Equal(String("  1  "), SUT1)
		asserts.Equal(String("  11 "), SUT2)
		asserts.Equal(String(" 121 "), SUT3)
	})

	t.Run("PadLeft", func(t *testing.T) {
		SUT1 := PadLeft(5, '.', "1")
		SUT2 := PadLeft(5, '.', "11")
		SUT3 := PadLeft(5, '.', "121")

		asserts.Equal(String("....1"), SUT1)
		asserts.Equal(String("...11"), SUT2)
		asserts.Equal(String("..121"), SUT3)
	})

	t.Run("PadRight", func(t *testing.T) {
		SUT1 := PadRight(5, '.', "1")
		SUT2 := PadRight(5, '.', "11")
		SUT3 := PadRight(5, '.', "121")

		asserts.Equal(String("1...."), SUT1)
		asserts.Equal(String("11..."), SUT2)
		asserts.Equal(String("121.."), SUT3)
	})

	t.Run("Trim", func(t *testing.T) {
		SUT := Trim("  hats  \n")

		asserts.Equal(String("hats"), SUT)
	})

	t.Run("TrimLeft", func(t *testing.T) {
		SUT := TrimLeft("  hats  \n")

		asserts.Equal(String("hats  \n"), SUT)
	})

	t.Run("TrimRight", func(t *testing.T) {
		SUT := TrimRight("  hats  \n")

		asserts.Equal(String("  hats"), SUT)
	})

}

func TestHigherOrderFinctions(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Map", func(t *testing.T) {
		SUT := Map(func(c char.Char) char.Char {
			if c == '/' {
				return '.'
			} else {
				return c
			}
		}, "a/b/c")

		asserts.Equal(String("a.b.c"), SUT)
	})

	t.Run("Filter", func(t *testing.T) {
		SUT := Filter(char.IsDigit, "R2-D2")

		asserts.Equal(String("22"), SUT)
	})

	t.Run("Foldl", func(t *testing.T) {
		SUT := Foldl(Cons, "", "time")

		asserts.Equal(String("emit"), SUT)
	})

	t.Run("Foldr", func(t *testing.T) {
		SUT := Foldr(Cons, "", "time")

		asserts.Equal(String("time"), SUT)
	})

	t.Run("Any", func(t *testing.T) {
		SUT1 := Any(char.IsDigit, "90210")
		SUT2 := Any(char.IsDigit, "R2-D2")
		SUT3 := Any(char.IsDigit, "heart")

		asserts.True(SUT1)
		asserts.True(SUT2)
		asserts.False(SUT3)
	})

	t.Run("All", func(t *testing.T) {
		SUT1 := All(char.IsDigit, "90210")
		SUT2 := All(char.IsDigit, "R2-D2")
		SUT3 := All(char.IsDigit, "heart")

		asserts.True(SUT1)
		asserts.False(SUT2)
		asserts.False(SUT3)
	})

}

func TestBasicsComparisons(t *testing.T) {
	asserts := assert.New(t)

	asserts.Equal(String("xyz"), basics.Max(String("abc"), String("xyz")))
	asserts.Equal(String("abc"), basics.Min(String("abc"), String("xyz")))
}
