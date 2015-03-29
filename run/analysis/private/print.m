function print(observed, predicted)
  [m1, v1] = report('Observed', observed);
  fprintf('\n');

  [m2, v2] = report('Predicted', predicted);
  fprintf('\n');

  fprintf('Error\n');
  fprintf('  Expectation: %.4e (%.2f%%)\n', ...
    abs(m1 - m2), 100 * abs((m1 - m2) / m1));
  fprintf('  Variance:    %.4e (%.2f%%)\n', ...
    abs(v1 - v2), 100 * abs((v1 - v2) / v1));
end

function [m, v] = report(title, data)
  m = mean(data);
  v = var(data);

  fprintf('%s\n', title);
  fprintf('  Expectation: %.4e\n', m);
  fprintf('  Variance:    %.4e\n', v);
end
