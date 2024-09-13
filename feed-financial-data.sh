#!/usr/bin/env bash

bq --format=csv query --use_legacy_sql=false "
SELECT
  partner_name,
  partner_id,
  ROUND(SUM(brl_amount), 2) AS volume,
  COUNT(*) AS operations,
  FORMAT_DATE('%F', created_at) created_at,
  CASE
    WHEN offer_type IN ('MESA', 'CORRESPONDENTES') THEN offer_type
    ELSE 'API'
END
  AS offer,
  CASE
    WHEN offer_type IN ('MESA', 'CORRESPONDENTES') THEN 'FX Core'
    ELSE 'DigitalFX (FXaaS)'
END
  AS product
FROM
  \`bexs-digitalfx.exchange.operations_normalized_latest\`
WHERE
  created_at >= DATE_TRUNC(CURRENT_DATE('America/Sao_Paulo'), month)
  AND status NOT IN ('ABORTED')
GROUP BY
  partner_name,
  partner_id,
  FORMAT_DATE('%F', created_at),
  CASE
    WHEN offer_type IN ('MESA', 'CORRESPONDENTES') THEN offer_type
    ELSE 'API'
END
  ,
  CASE
    WHEN offer_type IN ('MESA', 'CORRESPONDENTES') THEN 'FX Core'
    ELSE 'DigitalFX (FXaaS)'
END
ORDER BY
  5 DESC,
  1,
  6" > report.csv

./financial --source="report.csv"