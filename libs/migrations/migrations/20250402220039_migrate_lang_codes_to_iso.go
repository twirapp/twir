package migrations

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upMigrateLangCodesToIso, downMigrateLangCodesToIso)
}

func getIsoCode(lang int) string {
	// Map the enum values to ISO 639-1 codes
	isoMap := map[int]string{
		0:  "af", // Afrikaans
		1:  "sq", // Albanian
		2:  "ar", // Arabic
		3:  "hy", // Armenian
		4:  "az", // Azerbaijani
		5:  "eu", // Basque
		6:  "be", // Belarusian
		7:  "bn", // Bengali
		8:  "nb", // Bokmal
		9:  "bs", // Bosnian
		10: "bg", // Bulgarian
		11: "ca", // Catalan
		12: "zh", // Chinese
		13: "hr", // Croatian
		14: "cs", // Czech
		15: "da", // Danish
		16: "nl", // Dutch
		17: "en", // English
		18: "eo", // Esperanto
		19: "et", // Estonian
		20: "fi", // Finnish
		21: "fr", // French
		22: "lg", // Ganda
		23: "ka", // Georgian
		24: "de", // German
		25: "el", // Greek
		26: "gu", // Gujarati
		27: "he", // Hebrew
		28: "hi", // Hindi
		29: "hu", // Hungarian
		30: "is", // Icelandic
		31: "id", // Indonesian
		32: "ga", // Irish
		33: "it", // Italian
		34: "ja", // Japanese
		35: "kk", // Kazakh
		36: "ko", // Korean
		37: "la", // Latin
		38: "lv", // Latvian
		39: "lt", // Lithuanian
		40: "mk", // Macedonian
		41: "ms", // Malay
		42: "mi", // Maori
		43: "mr", // Marathi
		44: "mn", // Mongolian
		45: "nn", // Nynorsk
		46: "fa", // Persian
		47: "pl", // Polish
		48: "pt", // Portuguese
		49: "pa", // Punjabi
		50: "ro", // Romanian
		51: "ru", // Russian
		52: "sr", // Serbian
		53: "sn", // Shona
		54: "sk", // Slovak
		55: "sl", // Slovene
		56: "so", // Somali
		57: "st", // Sotho
		58: "es", // Spanish
		59: "sw", // Swahili
		60: "sv", // Swedish
		61: "tl", // Tagalog
		62: "ta", // Tamil
		63: "te", // Telugu
		64: "th", // Thai
		65: "ts", // Tsonga
		66: "tn", // Tswana
		67: "tr", // Turkish
		68: "uk", // Ukrainian
		69: "ur", // Urdu
		70: "vi", // Vietnamese
		71: "cy", // Welsh
		72: "xh", // Xhosa
		73: "yo", // Yoruba
		74: "zu", // Zulu
		75: "",   // Unknown
	}

	if code, exists := isoMap[lang]; exists {
		return code
	}
	return ""
}

func upMigrateLangCodesToIso(ctx context.Context, tx *sql.Tx) error {
	// Add new column for ISO codes
	_, err := tx.ExecContext(
		ctx, `
		ALTER TABLE channels_moderation_settings
		ADD COLUMN denied_chat_languages_iso text[] DEFAULT '{}'::text[]
	`,
	)
	if err != nil {
		return err
	}

	// Get all IDs first
	rows, err := tx.QueryContext(
		ctx, `
		SELECT id, denied_chat_languages::text[]
		FROM channels_moderation_settings
		WHERE denied_chat_languages IS NOT NULL
	`,
	)
	if err != nil {
		return err
	}

	type record struct {
		id          string
		deniedLangs []int64
	}

	var records []record
	for rows.Next() {
		var id string
		var arrayStr string
		if err := rows.Scan(&id, &arrayStr); err != nil {
			rows.Close()
			return err
		}

		// Parse the array string into []int64
		var deniedLangs []int64
		if arrayStr != "{}" && arrayStr != "" {
			// Remove braces and split
			numStr := arrayStr[1 : len(arrayStr)-1]
			if numStr != "" {
				nums := pq.Int64Array{}
				if err := nums.Scan(arrayStr); err != nil {
					rows.Close()
					return err
				}
				deniedLangs = []int64(nums)
			}
		}

		records = append(records, record{id: id, deniedLangs: deniedLangs})
	}
	rows.Close()

	// Process each record
	for _, rec := range records {
		var isoCodes []string
		for _, lang := range rec.deniedLangs {
			if isoCode := getIsoCode(int(lang)); isoCode != "" {
				isoCodes = append(isoCodes, isoCode)
			}
		}

		// Update the row with ISO codes
		_, err = tx.ExecContext(
			ctx, `
			UPDATE channels_moderation_settings
			SET denied_chat_languages_iso = $1
			WHERE id = $2
		`, pq.Array(isoCodes), rec.id,
		)
		if err != nil {
			return err
		}
	}

	// Drop the old column
	_, err = tx.ExecContext(
		ctx, `
		ALTER TABLE channels_moderation_settings
		DROP COLUMN denied_chat_languages
	`,
	)
	if err != nil {
		return err
	}

	// Rename the new column to the original name
	_, err = tx.ExecContext(
		ctx, `
		ALTER TABLE channels_moderation_settings
		RENAME COLUMN denied_chat_languages_iso TO denied_chat_languages
	`,
	)
	return err
}

func downMigrateLangCodesToIso(ctx context.Context, tx *sql.Tx) error {
	// Since we're losing information in the up migration (converting to ISO),
	// we'll just create an empty array for the down migration
	_, err := tx.ExecContext(
		ctx, `
		ALTER TABLE channels_moderation_settings
		ADD COLUMN denied_chat_languages_old int[] DEFAULT '{}'::int[]
	`,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx, `
		ALTER TABLE channels_moderation_settings
		DROP COLUMN denied_chat_languages
	`,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx, `
		ALTER TABLE channels_moderation_settings
		RENAME COLUMN denied_chat_languages_old TO denied_chat_languages
	`,
	)
	return err
}
