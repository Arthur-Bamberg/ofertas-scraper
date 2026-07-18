# Match-or-create Produto and Marca by normalized exact label

Produto and Marca identities are resolved by normalized exact string match (trim, lowercase, collapse internal whitespace): reuse the existing id when the normalized label matches, otherwise create. No fuzzy matching or human review queue in the MVP. The same policy applies to both Produto and Marca. Fuzzy match is a future upgrade if duplicates become costly.
